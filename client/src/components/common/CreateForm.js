import React from "react";

import {
  Form,
  TextInput,
  Button,
  Dropdown,
  TileGroup,
  RadioTile,
  Row,
  Grid,
  Column,
  FormLabel,
  Tooltip,
} from "carbon-components-react";

import geos from "../data/geo";

import styles from "./CreateForm.module.css";

import "./CreateForm.css";

const Spacer = ({ height }) => {
  return <div style={{ marginTop: height }} />;
};

const grab = async (url, options) => {
  const response = await fetch(url, options);
  if (response.status !== 200) {
    throw Error();
  }
  const data = await response.json();
  return data;
};

const CreateForm = () => {
  const [kubernetesSelected, setKubernetesSelected] = React.useState(true);
  const [openshiftSelected, setOpenshiftSelected] = React.useState(false);
  const [kuberntesVersions, setKubernetesVersions] = React.useState([]);
  const [openshiftVersions, setOpenshiftVersions] = React.useState([]);
  const [workerZones, setWorkerZones] = React.useState([]);
  const [privateVlans, setPrivateVlans] = React.useState([]);
  const [publicVlans, setPublicVlans] = React.useState([]);
  const [selectedKubernetes, setSelectedKuberetes] = React.useState(null);
  const [selectedOpenshift, setSelectedOpenshift] = React.useState(null);
  const [selectedRegion, setSelectedRegion] = React.useState(null);
  const [selectedWorkerZone, setSelectedWorkerZone] = React.useState(null);
  const [selectedPrivateVlan, setSelecetedPrivateVlan] = React.useState(null);
  const [selectedPublicVlan, setSelecetedPublicVlan] = React.useState(null);

  React.useEffect(() => {
    const loadVersions = async() => {
      try{
        const versions = await grab("/api/v1/clusters/versions");
        if (versions) {
          setKubernetesVersions(versions.kubernetes);
          setOpenshiftVersions(versions.openshift);
        }
      } catch (e) {
        console.log(e);
      }
    }
    loadVersions();
  }, []);

  const toggleRadio = () => {
    setKubernetesSelected(!kubernetesSelected);
    setOpenshiftSelected(!openshiftSelected);
  };

  const getWorkerZones = async (geo) => {
    console.log(geo);
    try {
      const locations = await grab(`/api/v1/clusters/${geo}/locations`, {
        Method: "GET",
      });
      if (locations && locations.length > 0) {
        setWorkerZones(locations);
      }
    } catch (e) {
      console.log(e);
    }
  };

  const getVlans = async (datacenter) => {
    try {
      const vlans = await grab(`/api/v1/clusters/${datacenter}/vlans`);
      const privateVlans = vlans.filter((vlan) => vlan.type === "private");
      if (privateVlans && privateVlans.length > 0) {
        setPrivateVlans(privateVlans);
        setSelecetedPrivateVlan(privateVlans[0]);
      }
      const publicVlans = vlans.filter((vlan) => vlan.type === "public");
      if (publicVlans && publicVlans.length > 0) {
        setPublicVlans(publicVlans);
        setSelecetedPublicVlan(publicVlans[0]);
      }
    } catch (e) {
      console.log(e);
    }
  };



  const getVersionString = (versions, version) => {
    const index = versions.indexOf(version);
    let substring = "stable";
    if (index === versions.length - 1) {
      substring = "latest";
    } else if (version.default) {
      substring = "stable, default"
    }

    return `${version.major}.${version.minor}.${version.patch} (${substring})`
  }

  const getKuberntesVersionString = (version) => {
    return getVersionString(kuberntesVersions, version);
  }

  const getOpenshiftVersionString = (version) => {
    return getVersionString(openshiftVersions, version);
  }

  const getVlanString = (vlan) => {
    return `${vlan.id}-${vlan.properties.vlan_number}-${vlan.properties.primary_router}`;
  };

  const onGeoSelected = (geo) => {
    getWorkerZones(geo.id);
    setSelectedRegion(geo);
    setSelectedWorkerZone(null);
    setPrivateVlans([]);
    setSelecetedPrivateVlan(null);
    setPublicVlans([]);
    setSelecetedPublicVlan(null);
  };

  const onWorkerZoneSelected = (worker) => {
    setSelectedWorkerZone(worker);
    getVlans(worker.id);
  };

  return (
    <Form>
      <Grid>
        <FormLabel>
          <Tooltip triggerText="Cluster type and version">
            The container platform type and version for the cluster. Choose
            Kubernetes for a native Kubernetes experience on Ubuntu, or
            OpenShift to deliver Kubernetes apps on Red Hat Enterprise Linux.
          </Tooltip>
        </FormLabel>
        <Row>
          <Column md={4} lg={3}>
            <RadioTile
              value="k8s"
              onClick={toggleRadio}
              checked={kubernetesSelected}
              className={styles.radio_tile}
            >
              <div className="radio-tile-content">
                <img
                  className="radio-tile-image"
                  src="https://i.ibb.co/cDqxKBd/download.png"
                  height={100}
                  width={100}
                />
                <p>Kubernetes</p>
              </div>
              <Dropdown
                id="kubernetes_version"
                className="create-page-dropdown"
                disabled={!kubernetesSelected}
                label="Select Version"
                items={kuberntesVersions}
                onChange={({selectedItem}) => setSelectedKuberetes(selectedItem)}
                selectedItem={selectedKubernetes}
                itemToString={version => getKuberntesVersionString(version)}
              />
            </RadioTile>
          </Column>
          <Column md={4} lg={3}>
            <RadioTile
              value="openshift"
              onClick={toggleRadio}
              checked={openshiftSelected}
              className={styles.radio_tile}
            >
              <div className="radio-tile-content">
                <img
                  className="radio-tile-image"
                  src="https://i.ibb.co/0fFQCD2/openshift.png"
                  height={100}
                  width={100}
                />
                <p>OpenShift</p>
              </div>
              <Dropdown
                id="openshift_version"
                className="create-page-dropdown"
                disabled={!openshiftSelected}
                label="Select Version"
                className={styles.dropdown}
                items={openshiftVersions}
                onChange={({selectedItem}) => setSelectedOpenshift(selectedItem)}
                selectedItem={selectedOpenshift}
                itemToString={version => getOpenshiftVersionString(version)}
              />
            </RadioTile>
          </Column>
        </Row>
        <Spacer height="16px" />

        <Spacer height="16px" />
        <Row>
          <Column>
            <h2>Location</h2>
          </Column>
        </Row>
        <Spacer height="16px" />

        <FormLabel>Geography</FormLabel>

        <Row>
          <Column lg={6}>
            <Dropdown
              id="geo_select"
              className="create-page-dropdown"
              label="Select geo"
              items={geos}
              selectedItem={selectedRegion}
              itemToString={(geo) => (geo ? geo.display_name : "")}
              onChange={({ selectedItem }) => onGeoSelected(selectedItem)}
            />
          </Column>
        </Row>
        <Spacer height="16px" />

        <Row>
          <Column lg={6}>
            <FormLabel>
              <Tooltip triggerText="Worker zone">
                The data center where your worker pool will be located.
              </Tooltip>
            </FormLabel>
            <Dropdown
              id="worker_zone_select"
              className="create-page-dropdown"
              label="Select worker zone"
              itemToString={(zone) => (zone ? zone.id : "")}
              items={workerZones}
              selectedItem={selectedWorkerZone}
              disabled={workerZones.length <= 0}
              onChange={({ selectedItem }) =>
                onWorkerZoneSelected(selectedItem)
              }
            />
          </Column>

          <Column md={4} lg={3}>
            <FormLabel>
              <Tooltip triggerText="Public vlan">
                Allow your worker nodes to securely communicate to the
                IBM-managed master through a virtual network. To expose your
                apps to the public, configure external networking.
              </Tooltip>
            </FormLabel>
            <Dropdown
              id="public_vlan"
              className="create-page-dropdown"
              label="Select public vlan"
              disabled={publicVlans.length <= 0}
              items={publicVlans}
              itemToString={(vlan) => (vlan ? getVlanString(vlan) : "")}
              selectedItem={selectedPublicVlan}
              onChange={({ selectedItem }) =>
                setSelecetedPublicVlan(selectedItem)
              }
            />
          </Column>

          <Column md={4} lg={3}>
            <FormLabel>
              <Tooltip triggerText="Private vlan">
                Virtual network that allows private communication between worker
                nodes in this cluster.
              </Tooltip>
            </FormLabel>
            <Dropdown
              id="private_vlan"
              className="create-page-dropdown"
              label="Select private vlan"
              disabled={privateVlans.length <= 0}
              items={privateVlans}
              itemToString={(vlan) => (vlan ? getVlanString(vlan) : "")}
              selectedItem={selectedPrivateVlan}
              onChange={({ selectedItem }) =>
                setSelecetedPrivateVlan(selectedItem)
              }
            />
          </Column>
        </Row>
        <Spacer height="16px" />

        <Spacer height="16px" />
        <Row>
          <Column>
            <h2>Cluster Metadata</h2>
          </Column>
        </Row>
        <Spacer height="16px" />

        <Row>
          <Column md={4} lg={3}>
            <FormLabel>
              <Tooltip triggerText="Cluster name prefix">
                Cluster name will be generated as prefix-001
              </Tooltip>
            </FormLabel>
            <TextInput labelText="" id="cluster_name" placeholder="mycluster" />
          </Column>

          <Column md={4} lg={3}>
            <FormLabel>Cluster count</FormLabel>
            <TextInput labelText="" id="cluster_count" placeholder="20" />
          </Column>

          <Column lg={3}>
            <FormLabel>Resouce group</FormLabel>
            <Dropdown
              id="resource_group"
              className="create-page-dropdown"
              label="Select resource group"
              items={["1", "2", "3"]}
            />
          </Column>
        </Row>
        <Spacer height="16px" />

        <Row>
          <Column>
            <FormLabel>
              <Tooltip triggerText="Tags">
                Use tags that would uniquely identify this set of clusters.
              </Tooltip>
            </FormLabel>
            <TextInput
              labelText=""
              id="tag_text"
              className="tag-text-input"
              placeholder="tag1, tag2, tag3"
            />
          </Column>
        </Row>
        <Spacer height="16px" />

        <Spacer height="16px" />
        <Row>
          <Column>
            <h2>Default worker pool</h2>
          </Column>
        </Row>
        <Spacer height="16px" />

        <Row>
          <Column md={4} lg={3}>
            <FormLabel>Worker nodes</FormLabel>
            <TextInput labelText="" id="worker_nodes" placeholder="1" />
          </Column>
        </Row>
        <Spacer height="16px" />

        <Row>
          <Column lg={6}>
            <FormLabel>
              <Tooltip triggerText="Flavor">
                The amount of memory, cpu, and disk space allocated to each
                worker node.
              </Tooltip>
            </FormLabel>
            <Dropdown
              id="machine_flavor"
              className="create-page-dropdown machine-flavor"
              label="Select Flavor"
              items={["1", "2", "3"]}
            />
          </Column>
        </Row>
        <Spacer height="16px" />
      </Grid>
    </Form>
  );
};

export default CreateForm;
