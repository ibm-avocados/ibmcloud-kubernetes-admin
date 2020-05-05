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
  const [workerZones, setWorkerZones] = React.useState([]);
  const [privateVlans, setPrivateVlans] = React.useState([]);
  const [publicVlans, setPublicVlans] = React.useState([]);

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
      setWorkerZones(locations);
    } catch (e) {
      console.log(e);
    }
  };

  const getVlans = async (datacenter) => {
    try {
      const vlans = await grab(`/api/v1/clusters/${datacenter}/vlans`);
      setPrivateVlans(vlans.filter((vlan) => vlan.type === "private"));
      setPublicVlans(vlans.filter((vlan) => vlan.type === "public"));
    } catch (e) {
      console.log(e);
    }
  };

  const getVlanString = (vlan) => {
    return `${vlan.id}-${vlan.properties.vlan_number}-${vlan.properties.primary_router}`;
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
                className="create-page-dropdown"
                disabled={!kubernetesSelected}
                label="Select Version"
                items={["1", "2", "3"]}
                selectedItem={null}
              />
            </RadioTile>
          </Column>
          <Column md={4} lg={3}>
            <RadioTile
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
                className="create-page-dropdown"
                disabled={!openshiftSelected}
                label="Select Version"
                className={styles.dropdown}
                items={["1", "2", "3"]}
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
              className="create-page-dropdown"
              label="Select geo"
              items={geos}
              itemToString={(geo) => (geo ? geo.display_name : "")}
              onChange={({ selectedItem }) => getWorkerZones(selectedItem.id)}
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
              className="create-page-dropdown"
              label="Select worker zone"
              itemToString={(zone) => (zone ? zone.id : "")}
              items={workerZones}
              disabled={workerZones.length <= 0}
              onChange={({ selectedItem }) => getVlans(selectedItem.id)}
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
              className="create-page-dropdown"
              label="Select public vlan"
              items={publicVlans}
              itemToString={(vlan) => getVlanString(vlan)}
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
              className="create-page-dropdown"
              label="Select private vlan"
              items={privateVlans}
              itemToString={(vlan) => getVlanString(vlan)}
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
            <TextInput placeholder="mycluster" />
          </Column>

          <Column md={4} lg={3}>
            <FormLabel>Cluster count</FormLabel>
            <TextInput placeholder="20" />
          </Column>

          <Column lg={3}>
            <FormLabel>Resouce group</FormLabel>
            <Dropdown
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
            <TextInput placeholder="1" />
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
