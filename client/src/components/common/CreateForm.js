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

import styles from "./CreateForm.module.css";

import "./CreateForm.css";

const Spacer = ({ height }) => {
  return <div style={{ marginTop: height }} />;
};

const CreateForm = () => {
  const [kubernetesSelected, setKubernetesSelected] = React.useState(true);
  const [openshiftSelected, setOpenshiftSelected] = React.useState(false);

  const toggleRadio = () => {
    setKubernetesSelected(!kubernetesSelected);
    setOpenshiftSelected(!openshiftSelected);
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
                disabled={!kubernetesSelected}
                label="Select Version"
                items={["1", "2", "3"]}
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
            <Dropdown label="Select geo" items={["1", "2", "3"]} />
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
            <Dropdown label="Select worker zone" items={["1", "2", "3"]} />
          </Column>

          <Column md={4} lg={3}>
            <FormLabel>
              <Tooltip triggerText="Public vlan">
                Allow your worker nodes to securely communicate to the
                IBM-managed master through a virtual network. To expose your
                apps to the public, configure external networking.
              </Tooltip>
            </FormLabel>
            <Dropdown label="Select public vlan" items={["1", "2", "3"]} />
          </Column>

          <Column md={4} lg={3}>
            <FormLabel>
              <Tooltip triggerText="Private vlan">
                Virtual network that allows private communication between worker
                nodes in this cluster.
              </Tooltip>
            </FormLabel>
            <Dropdown label="Select private vlan" items={["1", "2", "3"]} />
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
            <Dropdown label="Select resource group" items={["1", "2", "3"]} />
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
            <TextInput className="tag-text-input" placeholder="tag1, tag2, tag3"/>
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
            <FormLabel>
              Worker nodes
            </FormLabel>
            <TextInput placeholder="1"/>
          </Column>
        </Row>
        <Spacer height="16px" />

      </Grid>
    </Form>
  );
};

export default CreateForm;
