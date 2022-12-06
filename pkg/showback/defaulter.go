package showback

// +kubebuilder:webhook:path=/showback-defaulter-v1,mutating=true,failurePolicy=ignore,groups=core;apps;batch,resources=cronjobs;jobs;daemonsets;deployments;statefulsets,verbs=create;update,versions=*,name=showbackboss.shopify.io,sideEffects=None,admissionReviewVersions=v1

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"lmunro-at-shopify/test-defaulting-controller/pkg/utils"

	"k8s.io/apimachinery/pkg/runtime"
	informersV1 "k8s.io/client-go/informers/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	//+kubebuilder:scaffold:imports
)

var (
	logger  = ctrl.Log.WithName("defaulter")
	cluster = os.Getenv("CLUSTER")
)

const (
	bossName               = "showbackboss"
	showbackIDLabel        = "showback_id"
	showbackTypeLabel      = "showback_type"
	showbackStewardLabel   = "showback_steward_ref"
	showbackCostOwnerLabel = "showback_cost_owner_ref"
)

type Handler struct {
	Client            client.Client
	decoder           *admission.Decoder
	NamespaceInformer *informersV1.NamespaceInformer
}

func (h *Handler) addDefaultLabels(obj runtime.Object) (runtime.Object, error) {
	nsname, err := utils.RuntimeObjectNamespace(obj)
	if err != nil {
		return nil, err
	}
	objName, err := utils.RuntimeObjectName(obj)
	if err != nil {
		return nil, err
	}
	objKind, err := utils.RuntimeObjectKind(obj)
	if err != nil {
		return nil, err
	}

	namespace, err := (*h.NamespaceInformer).Lister().Get(nsname)
	if err != nil {
		return nil, err
	}
	nsLabels := namespace.ObjectMeta.GetLabels()
	defaultSteward, defaultStewardFound := nsLabels[showbackStewardLabel]
	defaultCostOwner, defaultCostOwnerFound := nsLabels[showbackCostOwnerLabel]

	objLabels, err := utils.RuntimeObjectLabels(obj)
	if err != nil {
		return nil, err
	}
	if _, found := objLabels[showbackIDLabel]; !found {
		obj, err = utils.RuntimeObjectSetLabels(obj, showbackIDLabel, objName)
		if err != nil {
			return nil, err
		}
	}
	if _, found := objLabels[showbackTypeLabel]; !found {
		obj, err = utils.RuntimeObjectSetLabels(obj, showbackTypeLabel, objKind)
		if err != nil {
			return nil, err
		}
	}

	podTemplate, err := utils.RuntimePodTemplate(obj)
	if err != nil {
		return nil, err
	}
	utils.SetPodSpecLabels(podTemplate, showbackIDLabel, objName)
	utils.SetPodSpecLabels(podTemplate, showbackTypeLabel, objKind)

	if defaultStewardFound {
		if _, found := objLabels[showbackStewardLabel]; !found {
			obj, err = utils.RuntimeObjectSetLabels(obj, showbackStewardLabel, defaultSteward)
			if err != nil {
				return nil, err
			}
		}
		utils.SetPodSpecLabels(podTemplate, showbackStewardLabel, defaultSteward)
	}

	if defaultCostOwnerFound {
		if _, found := objLabels[showbackCostOwnerLabel]; !found {
			obj, err = utils.RuntimeObjectSetLabels(obj, showbackCostOwnerLabel, defaultCostOwner)
			if err != nil {
				return nil, err
			}
		}
		utils.SetPodSpecLabels(podTemplate, showbackCostOwnerLabel, defaultCostOwner)
	}

	return obj, err
}

func (h *Handler) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger.Info("Handling Resource")
	obj, err := utils.DecodeFromGVK(req)
	if err != nil {
		logger.Error(err, "error decoding object")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	marshaledRequest, err := json.Marshal(req)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	logger.Info(string(marshaledRequest))

	obj, err = h.addDefaultLabels(obj)

	marshaled, err := json.Marshal(obj)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}

// InjectDecoder injects the decoder.
func (h *Handler) InjectDecoder(d *admission.Decoder) error {
	h.decoder = d
	return nil
}
