import React, { useCallback, useState } from 'react';
import {
  DataTable,
  DataTableSkeleton,
  Button,
  TableExpandRow,
  Tag,
  StructuredListWrapper,
  StructuredListHead,
  StructuredListBody,
  StructuredListRow,
  StructuredListCell,
  TextInput,
  SkeletonText,
  TagSkeleton,
  StructuredListSkeleton,
} from 'carbon-components-react';
import {
  Delete16 as Delete,
  TagGroup16 as TagGroup,
  Reset16 as Reset,
  Money16 as Money,
  VirtualMachine16 as VM,
} from '@carbon/icons-react';

import headers from '../../common/data/headers';

import './Cluster.css';
import useClusters from './useClusters';

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
  TableBatchAction,
} = DataTable;

const CustomExpandedRow = ({ name, dateCreated, workers, workersLoading, ownerEmail }) => (
  <>
    <h1>
      Cluster Name:
      {name}
    </h1>
    <h5>
      Date Created:
      {dateCreated}
    </h5>
    <h5>Owner Email: {ownerEmail}</h5>
    {workers ? <h3>Workers</h3> : <></>}
    {workers ? (
      <WorkerDetails workers={workers} />
    ) : workersLoading ? (
      <div style={{ width: '500px' }}>
        <StructuredListSkeleton rowCount={3} />
      </div>
    ) : (
      <div>Click on Load Workers to load worker details</div>
    )}
  </>
);

const WorkerDetails = ({ workers }) => (
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
      {workers.map((worker) => {
        const {
          id,
          state,
          machineType,
          privateVlan,
          publicVlan,
          status,
        } = worker;
        return (
          <StructuredListRow key={id}>
            <StructuredListCell noWrap>{state}</StructuredListCell>
            <StructuredListCell noWrap>{status}</StructuredListCell>
            <StructuredListCell>{publicVlan}</StructuredListCell>
            <StructuredListCell>{privateVlan}</StructuredListCell>
            <StructuredListCell>{machineType}</StructuredListCell>
          </StructuredListRow>
        );
      })}
    </StructuredListBody>
  </StructuredListWrapper>
);

const Clusters = ({ accountID }) => {
  const [
    clusters,
    { deleteClusters, deleteTag, setTag, reload, getBilling, getWorkers },
  ] = useClusters(accountID);

  console.log(clusters);

  const [tagText, setTagText] = useState('');
  const [billingLoading, setBillingLoading] = useState(false);
  const [workersLoading, setWorkersLoading] = useState(false);

  const onBillingClicked = useCallback(
    (data) => {
      setBillingLoading(true);
      getBilling(data);
    },
    [getBilling]
  );

  const onGetWorkersClicked = useCallback((data) => {
    setWorkersLoading(true);
    getWorkers(data);
  });

  const onSetTagClicked = useCallback(
    (clusters, tagText) => {
      setTagText('');
      setTag(clusters, tagText);
    },
    [setTag]
  );

  const CustomCell = ({ cell, crn, id }) => {
    const { info, value } = cell;
    switch (info.header) {
      case 'state':
        return (
          <span className="oneline">
            <span className={`status ${value}`} />
            {value}
          </span>
        );
      case 'masterKubeVersion':
        return (
          <span className="oneline">
            <img
              alt="logo"
              className="logo-image"
              src={
                value.includes('openshift')
                  ? 'https://i.ibb.co/tLktm91/os-icon.png'
                  : 'https://i.ibb.co/Hh2TzLH/k8s-icon.png'
              }
            />
            {value}
          </span>
        );
      case 'tags':
        return (
          <>
            {value ? (
              value.map((tag) => (
                <Tag
                  onClose={() => deleteTag(id, tag, crn)}
                  filter
                  key={tag}
                  type="blue"
                >
                  {tag}
                </Tag>
              ))
            ) : (
              <div>
                <TagSkeleton />
              </div>
            )}
          </>
        );
      case 'cost':
        if (value) {
          return <>${value}</>;
        }
        return (
          <>
            {billingLoading ? (
              <div style={{ width: '50px' }}>
                <SkeletonText />
              </div>
            ) : (
              '$'
            )}
          </>
        );
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
      getExpandHeaderProps,
    }) => (
      <TableContainer title="Clusters">
        <TableToolbar>
          {/* pass in `onInputChange` change here to make filtering work */}
          <TableBatchActions {...getBatchActionProps()}>
            <TableBatchAction
              tabIndex={getBatchActionProps().shouldShowBatchActions ? 0 : -1}
              renderIcon={Delete}
              onClick={() =>
                deleteClusters(selectedRows.map((r) => clusters.data[r.id]))
              }
            >
              Delete
            </TableBatchAction>
            <div className="tag-input">
              <TextInput
                id="tag-input"
                hideLabel
                value={tagText}
                onChange={(e) => setTagText(e.target.value.trim())}
                labelText="tag"
                placeholder="Tag"
              />
            </div>
            <Button
              renderIcon={TagGroup}
              iconDescription="Group Tag"
              hasIconOnly
              kind="primary"
              size="default"
              type="button"
              tooltipPosition="right"
              onClick={() =>
                onSetTagClicked(
                  selectedRows.map((r) => clusters.data[r.id]),
                  tagText
                )
              }
            />
          </TableBatchActions>
          <TableToolbarContent>
            <TableToolbarSearch
              tabIndex={getBatchActionProps().shouldShowBatchActions ? -1 : 0}
              onChange={onInputChange}
            />
            <Button
              renderIcon={VM}
              iconDescription="Get Workers for Clusters"
              hasIconOnly
              kind="secondary"
              size="field"
              type="button"
              tooltipPosition="right"
              onClick={() => onGetWorkersClicked(clusters.data)}
            />
            <Button
              renderIcon={Money}
              iconDescription="Get Billing for Clusters"
              hasIconOnly
              kind="tertiary"
              size="field"
              type="button"
              tooltipPosition="right"
              onClick={() => onBillingClicked(clusters.data)}
            />
            <Button onClick={reload} renderIcon={Reset}>
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
              <TableExpandHeader enableExpando {...getExpandHeaderProps()} />
              <TableSelectAll {...getSelectionProps()} />
              {headers.map((header, i) => (
                <TableHeader key={i} {...getHeaderProps({ header })}>
                  {header.header}
                </TableHeader>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row) => (
              <React.Fragment key={row.id}>
                <TableExpandRow {...getRowProps({ row })}>
                  <TableSelectRow {...getSelectionProps({ row })} />
                  {row.cells.map((cell) => (
                    <TableCell key={cell.id}>
                      <CustomCell
                        cell={cell}
                        crn={clusters.data[row.id].crn}
                        id={row.id}
                      />
                    </TableCell>
                  ))}
                </TableExpandRow>
                <TableExpandedRow colSpan={headers.length + 2}>
                  <CustomExpandedRow
                    ownerEmail={clusters.data[row.id].ownerEmail}
                    workersLoading={workersLoading}
                    name={clusters.data[row.id].name}
                    dateCreated={clusters.data[row.id].createdDate}
                    workers={clusters.data[row.id].workers}
                  />
                </TableExpandedRow>
              </React.Fragment>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    ),
    [
      clusters.data,
      deleteClusters,
      onBillingClicked,
      onSetTagClicked,
      reload,
      tagText,
    ]
  );

  if (clusters.isLoading) {
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
    <DataTable
      rows={Object.keys(clusters.data).map((id) => clusters.data[id])}
      headers={headers}
      render={render}
      isSortable
    />
  );
};

export default Clusters;
