// +kubebuilder:webhook:path=/showback-defaulter-v1,mutating=true,failurePolicy=ignore,groups=*,resources=pods;cronjobs;jobs;daemonsets;deployments;statefulsets;replicasets,verbs=create;update,versions=*,name=showback-boss,sideEffects=None,admissionReviewVersions=v1
package showback

import (
	"context"
	"encoding/json"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	informersV1 "k8s.io/client-go/informers/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	//+kubebuilder:scaffold:imports
)

const (
	showbackIDLabel        = "showback_id"
	showbackTypeLabel      = "showback_type"
	showbackStewardLabel   = "showback_steward_ref"
	showbackCostOwnerLabel = "showback_cost_owner_ref"
	deploymentTypeValue    = "gke_deployment"
)

type Defaulter struct {
	Client            client.Client
	decoder           *admission.Decoder
	NamespaceInformer *informersV1.NamespaceInformer
}

func (def *Defaulter) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger := ctrl.Log.WithName("defaulter")
	deployment := &appsv1.Deployment{}
	err := def.decoder.Decode(req, deployment)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	marshaledRequest, err := json.Marshal(req)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	logger.Info(string(marshaledRequest))

	namespace, err := (*def.NamespaceInformer).Lister().Get(deployment.GetNamespace())
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	nsLabels := namespace.ObjectMeta.GetLabels()
	deployLabels := &(deployment.Spec.Template.ObjectMeta.Labels)
	defaultSteward, defaultStewardFound := nsLabels[showbackStewardLabel]
	defaultCostOwner, defaultCostOwnerFound := nsLabels[showbackCostOwnerLabel]

	if _, found := (*deployLabels)[showbackIDLabel]; !found {
		(*deployLabels)[showbackIDLabel] = deployment.GetObjectMeta().GetName()
	}
	if _, found := (*deployLabels)[showbackTypeLabel]; !found {
		(*deployLabels)[showbackTypeLabel] = deploymentTypeValue
	}

	if _, found := (*deployLabels)[showbackStewardLabel]; defaultStewardFound && !found {
		(*deployLabels)[showbackStewardLabel] = defaultSteward
	}
	if _, found := (*deployLabels)[showbackCostOwnerLabel]; defaultCostOwnerFound && !found {
		(*deployLabels)[showbackCostOwnerLabel] = defaultCostOwner
	}

	marshaled, err := json.Marshal(deployment)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}

// InjectDecoder injects the decoder.
func (def *Defaulter) InjectDecoder(d *admission.Decoder) error {
	def.decoder = d
	return nil
}
