package main

import (
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// drop Update Events where the object's generation(spec) doesn't changed
var pdcDropBanalUpdate = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
	},
}

// ignore events happend in kube-system namespace
var pdcIgnoreSysNs = predicate.Funcs {
	CreateFunc: func(e event.CreateEvent) bool {
		return e.Object.GetNamespace() != "kube-system"
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		return e.Object.GetNamespace() != "kube-system"
	},
	UpdateFunc: func(e event.UpdateEvent) bool {
		return e.ObjectOld.GetNamespace() != "kube-system"
	},
	GenericFunc: func(e event.GenericEvent) bool {
		return e.Object.GetNamespace() != "kube-system"
	},
}

var anotPred = predicate.And(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.AnnotationChangedPredicate{}),pdcIgnoreSysNs)
