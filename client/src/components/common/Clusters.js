import React, { useCallback, useState } from "react";
import {
  DataTable,
  DataTableSkeleton,
  Button,
  TableExpandRow,
  Loading,
  Tag,
  StructuredListWrapper,
  StructuredListHead,
  StructuredListBody,
  StructuredListRow,
  StructuredListCell,
  TextInput
} from "carbon-components-react";
import {
  Delete16 as Delete,
  TagGroup16 as TagGroup,
  Reset16 as Reset
} from "@carbon/icons-react";

import headers from "../data/headers";

import "./Cluster.css";
import { useEffect } from "react";

const {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableBody,
  TableCell,
  TableHeader,
  TableExpandHeader,
  TableSelectRow,
  TableExpandedRow,
  TableSelectAll,
  TableToolbar,
  TableToolbarSearch,
  TableToolbarContent,
  TableBatchActions,
  TableBatchAction
} = DataTable;

// Takes an array of objects and tranforms it into a map of objects, with ID
// being the key and the object being the value.
// e.g.
// [{ id: 'a1', x: 'hello' }, { id: 'b2', x: 'world' }] =>
// {
//   a1: { id: 'a1', x: 'hello' },
//   b2: { id: 'b2', x: 'world' }
// }
const arrayToMap = arr =>
  arr.reduce((acc, cur) => ({ ...acc, [cur.id]: cur }), {});

// const mapToArray = data =>
//   Object.keys(data).map(val => data[val]);


const deleteCluster = cluster =>
  fetch("/api/v1/clusters", {
    method: "DELETE",
    body: JSON.stringify({
      id: cluster.id,
      resourceGroup: cluster.resourceGroup,
      deleteResources: true
    })
  });

const CustomExpandedRow = ({ name, dateCreated, workers }) => {
  return (
    <>
      <h1>Cluster Name: {name}</h1>
      <h5>Date Created: {dateCreated}</h5>
      {workers ? <h3>Workers</h3> : <></>}
      {workers ? <WorkerDetails workers={workers} /> : <></>}
    </>
  )
}

const WorkerDetails = ({ workers }) => {
  return (
    <StructuredListWrapper>
      <StructuredListHead>
        <StructuredListRow head>
          <StructuredListCell head>State</StructuredListCell>
          <StructuredListCell head>Status</StructuredListCell>
          <StructuredListCell head>Public Vlan</StructuredListCell>
          <StructuredListCell head>Private Vlan</StructuredListCell>
          <StructuredListCell head>Machine Type</StructuredListCell>
        </StructuredListRow>
      </StructuredListHead>
      <StructuredListBody>
        {workers.map(worker => {
          const { id, state, machineType, privateVlan, publicVlan, status } = worker;
          return (
            <StructuredListRow key={id}>
              <StructuredListCell noWrap>{state}</StructuredListCell>
              <StructuredListCell noWrap>{status}</StructuredListCell>
              <StructuredListCell>{publicVlan}</StructuredListCell>
              <StructuredListCell>{privateVlan}</StructuredListCell>
              <StructuredListCell>{machineType}</StructuredListCell>
            </StructuredListRow>
          )
        })}
      </StructuredListBody>
    </StructuredListWrapper>
  )
}

const Clusters = ({ accountID }) => {
  const [isLoadingClusters, setLoadingClusters] = useState(true);
  const [showLoading, setShowLoading] = useState(false);
  const [clusters, setClusters] = useState([]);
  const [tagText, setTagText] = useState("");
  // const [accountIDData, setAccountID] = useState(accountID);

  const loadClusters = useCallback(async () => {
    setLoadingClusters(true);
    const response = await fetch(`/api/v1/clusters/${accountID}`);
    if (response.status !== 200) {

    }
    const clusters = await response.json();
    console.log(clusters);
    setClusters(clusters);
    setLoadingClusters(false);
  }, [accountID]);

  useEffect(() => {
    loadClusters();
  }, [loadClusters]);






  const deleteClusters = useCallback(
    clusters => async () => {
      setShowLoading(true);
      const promises = clusters.map(cluster => deleteCluster(cluster));
      await Promise.all(promises);
      setShowLoading(false);
      loadClusters();
    },
    [loadClusters]
  );

  // const buttonClicked = rows => () => {
  //   console.log("slected rows", rows);
  // };

  const deleteTag = useCallback((tagName, crn) => async () => {
    let body = {
      tag_name: tagName,
      resources: [{ resource_id: crn }]
    };
    const response = await fetch("/api/v1/clusters/deletetag", {
      method: "POST",
      body: JSON.stringify(body)
    });

    if (response.status !== 200) {

    }  
    const result = await response.json();
    console.log(result);
    loadClusters();
  }, [loadClusters]);

  const setTag = useCallback(clusters => async () => {
    if (tagText === "") {
      return;
    }
    let resources = clusters.map(cluster => { return ({ resource_id: cluster.crn }) });
    let body = {
      tag_name: tagText,
      resources: resources
    };
    setShowLoading(true);
    const response = await fetch("/api/v1/clusters/settag", {
      method: "POST",
      body: JSON.stringify(body)
    });
    if (response.status !== 200) {

    }

    setShowLoading(false);
    loadClusters();
    setTagText("");
  }, [loadClusters, tagText]);

  const CustomCell = ({ cell, crn }) => {
    const { info, value } = cell;
    switch (info.header) {
      case "state":
        return (
          <span className="oneline">
            <span className={`status ${value}`}></span>
            {value}
          </span>
        );
      case "masterKubeVersion":
        return (
          <span className="oneline">
            <img
              alt="logo"
              className="logo-image"
              src={
                value.includes("openshift")
                  ? "https://cloud.ibm.com/kubernetes/img/openshift_logo-7825001afb.svg"
                  : "https://cloud.ibm.com/kubernetes/img/container-service-logo-7e87826329.svg"
              }
            />
            {value}
          </span>
        );
      case "tags":
        return (<>{value.map(tag => <Tag onClick={deleteTag(tag, crn)} filter key={tag} type='blue'>{tag}</Tag>)}</>)
      case "cost":
        return (<>${value}</>)
      default:
        return <>{value}</>;
    }
  };

  const render = useCallback(
    ({
      rows,
      headers,
      getHeaderProps,
      getRowProps,
      getBatchActionProps,
      getSelectionProps,
      selectedRows,
      onInputChange,
      getExpandHeaderProps
    }) => {
      const clusterMap = arrayToMap(clusters);
      return (
        <TableContainer title="Clusters">
          <TableToolbar>
            {/* pass in `onInputChange` change here to make filtering work */}
            <TableBatchActions {...getBatchActionProps()}>
              <TableBatchAction
                tabIndex={getBatchActionProps().shouldShowBatchActions ? 0 : -1}
                renderIcon={Delete}
                onClick={deleteClusters(
                  selectedRows.map(r => clusterMap[r.id])
                )}
              >
                Delete
              </TableBatchAction>
              <div className="tag-input">
                <TextInput
                  id="tag-input"
                  hideLabel
                  onChange={e => setTagText(e.target.value.trim())}
                  labelText="tag"
                  placeholder="Tag" />
              </div>
              <Button
                renderIcon={TagGroup}
                iconDescription="Group Tag"
                hasIconOnly kind="primary"
                size="default"
                type="button"
                tooltipPosition="bottom"
                onClick={setTag(
                  selectedRows.map(r => clusterMap[r.id])
                )} />
            </TableBatchActions>
            <TableToolbarContent>
              <TableToolbarSearch
                tabIndex={getBatchActionProps().shouldShowBatchActions ? -1 : 0}
                onChange={onInputChange}
              />
              <Button onClick={loadClusters} renderIcon={Reset}>
                Reload
              </Button>
            </TableToolbarContent>
            {/* <TableToolbarContent>
            <Button onClick={() => buttonClicked(selectedRows)}  kind="primary">
              Add new
            </Button>
          </TableToolbarContent> */}
          </TableToolbar>
          <Table>
            <TableHead>
              <TableRow>
                <TableExpandHeader
                  enableExpando={true}
                  {...getExpandHeaderProps()}
                />
                <TableSelectAll {...getSelectionProps()} />
                {headers.map(header => (
                  <TableHeader {...getHeaderProps({ header })}>
                    {header.header}
                  </TableHeader>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {rows.map(row => (
                <React.Fragment key={row.id}>
                  <TableExpandRow {...getRowProps({ row })}>
                    <TableSelectRow {...getSelectionProps({ row })} />
                    {row.cells.map(cell => (
                      <TableCell key={cell.id}>
                        <CustomCell cell={cell} crn={clusterMap[row.id].crn} />
                      </TableCell>
                    ))}
                  </TableExpandRow>
                  <TableExpandedRow colSpan={headers.length + 2}>
                    <CustomExpandedRow name={clusterMap[row.id].name}
                      dateCreated={clusterMap[row.id].createdDate}
                      workers={clusterMap[row.id].workers}
                    />

                  </TableExpandedRow>
                </React.Fragment>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      );
    },
    [clusters, deleteClusters, loadClusters, setTag]
  );

  if (isLoadingClusters) {
    return (
      <>
        <div className="bx--data-table-header">
          <h4>Clusters</h4>
        </div>
        <DataTableSkeleton
          columnCount={headers.length}
          headers={headers}
          rowCount={5}
          zebra
        />
      </>
    );
  }

  return (
    <>
      <DataTable
        rows={clusters}
        headers={headers}
        render={render}
        isSortable //={true}
      />
      <Loading active={showLoading} />)
    </>
  );
};

export default Clusters;
