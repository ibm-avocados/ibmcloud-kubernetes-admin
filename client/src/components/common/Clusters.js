import React, { useCallback, useState } from "react";
import { DataTable, DataTableSkeleton, Button } from "carbon-components-react";
import {
  Delete16 as Delete,
  Save16 as Save,
  Reset16 as Reset
} from "@carbon/icons-react";
import { getJSON } from "../../fetchUtil";

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
  TableSelectRow,
  TableSelectAll,
  TableToolbar,
  TableToolbarSearch,
  TableToolbarContent,
  TableBatchActions,
  TableBatchAction
} = DataTable;

// mmm so sexy and unreadable 🥵💦
const arrayToMap = arr =>
  arr.reduce((acc, cur) => ({ ...acc, [cur.id]: cur }), {});

const deleteCluster = cluster =>
  fetch("/api/v1/clusters", {
    method: "DELETE",
    body: JSON.stringify({
      id: cluster.id,
      resourceGroup: cluster.resourceGroup,
      deleteResources: true
    })
  });

const processHeader = header => {
  return header.header;
};

const process = cell => {
  let id = cell.id;
  let field = id.split(":")[1];
  let value = cell.value;
  if (field === "state") {
    if (value === "normal") {
      return (
        <span className="oneline">
          <span className="status normal"></span>
          {value}
        </span>
      );
    } else if (value === "warning") {
      return (
        <span className="oneline">
          <span className="status warning"></span>
          {value}
        </span>
      );
    } else if (value === "critical") {
      return (
        <span className="oneline">
          <span className="status critical"></span>
          {value}
        </span>
      );
    }
  } else if (field === "masterKubeVersion") {
    if (value.includes("openshift")) {
      return (
        <>
          <span className="oneline">
            <img
              alt="openshift logo"
              className="logo-image"
              src="https://cloud.ibm.com/kubernetes/img/openshift_logo-7825001afb.svg"
            />
            {value}
          </span>
        </>
      );
    }
    return (
      <>
        <span className="oneline">
          <img
            alt="iks logo"
            className="logo-image"
            src="https://cloud.ibm.com/kubernetes/img/container-service-logo-7e87826329.svg"
          />
          {value}
        </span>
      </>
    );
  }
  return <>{value}</>;
};

const Clusters = ({ accountChanged }) => {
  const [isLoadingClusters, setLoadingClusters] = useState(true);
  const [clusters, setClusters] = useState([]);

  const loadClusters = useCallback(async () => {
    setLoadingClusters(true);
    const clusters = await fetch("/api/v1/clusters").then(getJSON);
    setClusters(clusters);
    setLoadingClusters(false);
  }, []);

  useEffect(() => {
    loadClusters();
  }, [loadClusters]);

  const deleteClusters = useCallback(
    clusters => async () => {
      console.log(clusters);
      // const promises = clusters.map(cluster => deleteCluster(cluster));
      // await Promise.all(promises);
      loadClusters();
    },
    [loadClusters]
  );

  const render = useCallback(
    ({
      rows,
      headers,
      getHeaderProps,
      getBatchActionProps,
      getSelectionProps,
      selectedRows,
      onInputChange
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
              <TableBatchAction
                renderIcon={Save}
                onClick={() => alert("Do what now?")}
              >
                Save
              </TableBatchAction>
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
                <TableSelectAll {...getSelectionProps()} />
                {headers.map(header => (
                  <TableHeader {...getHeaderProps({ header })}>
                    {processHeader(header)}
                  </TableHeader>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {rows.map(row => {
                return (
                  <TableRow key={row.id}>
                    <TableSelectRow {...getSelectionProps({ row })} />
                    {row.cells.map(cell => (
                      <TableCell key={cell.id}>{process(cell)}</TableCell>
                    ))}
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
      );
    },
    [clusters, deleteClusters, loadClusters]
  );

  if (isLoadingClusters) {
    return (
      <>
        <div className="bx--data-table-header">
          <h4>Clusters</h4>
        </div>
        <DataTableSkeleton
          columnCount={headers.length}
          compact={false}
          headers={headers}
          rowCount={5}
          zebra={true}
        />
      </>
    );
  }

  return (
    <DataTable
      rows={clusters}
      headers={headers}
      render={render}
      isSortable={true}
    />
  );
};

export default Clusters;
