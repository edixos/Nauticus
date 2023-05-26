## Introduction

Welcome to the Nauticus Controller! This project provides a Kubernetes controller that allows you to easily
manage `Spaces`, a fully-managed Kubernetes namespace that integrates RBAC, network policies, resources, and quotas and
more. This getting started guide will walk you through the process of installing and configuring Nauticus for your
Kubernetes cluster, as well as give an overview of the features and functionality provided by the controller. We hope
that this guide will help you to quickly get up and running with Nauticus, and start taking advantage of the powerful
tooling it provides. Let's get started!

## Prerequisites

- A running Kubernetes cluster (version 1.17 or later) with kubectl command-line tool installed and configured on your
  machine.
- Access to the cluster with cluster-admin permissions.
- Familiarity with Kubernetes concepts such as namespaces, RBAC, and resources.

Before you begin, please ensure that you have a running Kubernetes cluster and you have the kubectl command-line tool
installed and configured on your machine. If you don't have a cluster, you can create one using the kind command. You
must also have cluster-admin permissions to be able to use Nauticus.

It is also recommended to have basic understanding of Kubernetes concepts such as namespaces, RBAC, and resources.

Once you have everything set up, you can proceed to the next step of the guide.

## Installation

In this section, we will explain how to install Nauticus using Helm and kubectl.

### Helm

1. Make sure that you have docker and kind installed if you want to create a kubernetes cluster within Docker
2. Start by creating a new Kubernetes cluster with kind command if you don't have one.

     ```bash title="Create a kind cluster"
     kind create cluster --image kindest/node:v1.27.2 --wait 5m --name nauticus
     ```

3. Make sure that you have Helm and kubectl installed on your machine. If you don't have them installed, you can find
   instructions on how to install them here: [helm](https://helm.sh/docs/intro/install/) and
   here [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

4. Use the following command to add the Nauticus Helm repository:

     ```bash  title="Add Edixos Helm Repository"
     helm repo add edixos https://edixos.github.io/charts
     ```
   
5. Use the following commands to install Nauticus:
      
     ```bash title="Install Naurticus with Helm"
     kubectl create namespace nauticus-system
     helm install nauticus --namespace nauticus-system edixos/nauticus
     ```
   
6. Use the following command to check the status of the Nauticus pods:
     
     ```bash  title="Get Nauticus Controller pod"
     kubectl get pods -n nauticus-system
     ```
   
7. Use the following command to verify that the Nauticus controller is running:
     ```bash  title="Get Nauticus Controller logs"
     kubectl logs -f <nauticus-controller-pod-name> -n nauticus-system
     ```

Now you can start using Nauticus to create and manage Spaces in your cluster.

### Kubectl

In this section, we will go through the steps to install Nauticus using a kubectl manifest file. This file contains all
the resources needed for Nauticus to function properly.

1. Apply the manifest file using kubectl:

     ```bash  title="Install Nauticus from all-in-one manifest file"
     kubectl apply -f https://raw.githubusercontent.com/edixos/nauticus/main/config/install.yaml
     ```

2. Use the following command to check the status of the Nauticus pods:
     ```bash  title="Get Nauticus Controller pod"
     kubectl get pods -n nauticus-system
     ```
3. Use the following command to verify that the Nauticus controller is running:
     ```bash title="Get Nauticus Controller logs"
     kubectl logs -f <nauticus-controller-pod-name> -n nauticus-system
     ```
   
## Basic Usage

To use Nauticus to create a new space, you can use the following kubectl command:

  ```bash title="Create a basic Space"
cat << EOF | kubectl apply -f -
apiVersion: nauticus.io/v1alpha1
kind: Space
metadata:
  name: my-space
EOF
  ```


The space will create a namespace and update the status of the `Space` by adding `status.NamespaceName` to it.
For more details on how to use space's features please refer to the Tutorial Section.

