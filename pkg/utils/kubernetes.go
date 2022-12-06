package utils

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	jsonSerializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func RuntimeObjectKind(obj runtime.Object) (string, error) {
	objUnstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return "", err
	}
	unstructured := &unstructured.Unstructured{
		Object: objUnstr,
	}
	return unstructured.GetKind(), nil
}

func RuntimeObjectName(obj runtime.Object) (string, error) {
	objUnstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return "", err
	}
	unstructured := &unstructured.Unstructured{
		Object: objUnstr,
	}
	return unstructured.GetName(), nil
}

func RuntimeObjectNamespace(obj runtime.Object) (string, error) {
	objUnstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return "", err
	}
	unstructured := &unstructured.Unstructured{
		Object: objUnstr,
	}
	return unstructured.GetNamespace(), nil
}

func RuntimeObjectOwnerReference(obj runtime.Object) ([]metav1.OwnerReference, error) {
	blankOR := make([]metav1.OwnerReference, 0)
	objUnstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return blankOR, err
	}
	unstructured := &unstructured.Unstructured{
		Object: objUnstr,
	}
	return unstructured.GetOwnerReferences(), nil
}

func RuntimeObjectSetLabels(obj runtime.Object, key, value string) (runtime.Object, error) {
	objUnstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return obj, err
	}
	unstructured := &unstructured.Unstructured{
		Object: objUnstr,
	}
	labels := unstructured.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels[key] = value
	unstructured.SetLabels(labels)
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.Object, obj)
	return obj, err
}

func RuntimeObjectLabels(obj runtime.Object) (map[string]string, error) {
	blankLabels := make(map[string]string)
	objUnstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return blankLabels, err
	}
	unstructured := &unstructured.Unstructured{
		Object: objUnstr,
	}
	return unstructured.GetLabels(), nil
}

// DecodeFromGVK decodes the object into a runtime.Object according to the GVK
// info that the request has for the object. Only use this function if you need
// to support more than one type of resource
func DecodeFromGVK(req admission.Request) (runtime.Object, error) {
	s := jsonSerializer.NewYAMLSerializer(jsonSerializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	gvk := schema.GroupVersionKind(req.AdmissionRequest.Kind)
	obj, _, err := s.Decode(req.AdmissionRequest.Object.Raw, &gvk, nil)
	return obj, err
}

func RuntimePodTemplate(obj runtime.Object) (*v1.PodTemplateSpec, error) {
	switch t := obj.(type) {
	case *batchv1beta1.CronJob:
		return &t.Spec.JobTemplate.Spec.Template, nil
	case *batchv1.CronJob:
		return &t.Spec.JobTemplate.Spec.Template, nil
	case *appsv1.DaemonSet:
		return &t.Spec.Template, nil
	case *extensionsv1beta1.DaemonSet:
		return &t.Spec.Template, nil
	case *appsv1beta2.DaemonSet:
		return &t.Spec.Template, nil
	case *extensionsv1beta1.Deployment:
		return &t.Spec.Template, nil
	case *appsv1.Deployment:
		return &t.Spec.Template, nil
	case *appsv1beta1.Deployment:
		return &t.Spec.Template, nil
	case *appsv1beta2.Deployment:
		return &t.Spec.Template, nil
	case *batchv1.Job:
		return &t.Spec.Template, nil
	case *appsv1.ReplicaSet:
		return &t.Spec.Template, nil
	case *extensionsv1beta1.ReplicaSet:
		return &t.Spec.Template, nil
	case *appsv1beta2.ReplicaSet:
		return &t.Spec.Template, nil
	case *corev1.ReplicationController:
		return t.Spec.Template, nil
	case *appsv1.StatefulSet:
		return &t.Spec.Template, nil
	case *appsv1beta1.StatefulSet:
		return &t.Spec.Template, nil
	case *appsv1beta2.StatefulSet:
		return &t.Spec.Template, nil
	default:
		return nil, fmt.Errorf("unknown kind: %v", t)
	}
}

func SetPodSpecLabels(template *v1.PodTemplateSpec, key, value string) {
	if key == "" || value == "" {
		return
	}
	if (*template).ObjectMeta.Labels == nil {
		(*template).ObjectMeta.Labels = make(map[string]string)
	}
	if _, found := (*template).ObjectMeta.Labels[key]; !found {
		(*template).ObjectMeta.Labels[key] = value
	}
}
