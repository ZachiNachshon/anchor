#!/usr/bin/env bash

function install() {
	local output=$(kubectl get deployments kubernetes-dashboard --no-headers 2>&1)

	if [[ $output != *"NotFound"* ]]; then
		echo "===> Dashboard already installed; skipping"
		return
	fi

	kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
	kubectl create serviceaccount -n kube-system kubernetes-dashboard
	kubectl create clusterrolebinding -n kube-system kubernetes-dashboard \
		--clusterrole cluster-admin \
		--serviceaccount kube-system:kubernetes-dashboard
}

function start() {
	echo && echo "===> Token:"
	kubectl -n kube-system describe secret $(
		kubectl -n kube-system get secret | awk '/^kubernetes-dashboard-token-/{print $1}'
	) | awk '$1=="token:"{print $2}'

	open "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy"
	echo && echo && kubectl proxy
}

function init() {
	export KUBECONFIG="$(kind get kubeconfig-path --name=kind)"
}

function main() {
	install
	start
}

init
main
