import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import RequestForm from './components/requestform';
import RequestTable from './components/requesttable';
import ClusterTable from './components/clustertable';

const useStyles = makeStyles((theme) => ({
  clusterpanel: {
    display: 'flex',
    flex: 1,
    flexDirection: 'column',
  },
  form: {
    padding: "10px",
  },
  tables: {
    display: 'flex',
    flex: 1,
    flexDirection: 'row',
    overflow: "hidden",
  },
  table: {
    flex: 1,
    overflow: 'auto',
    padding: "10px",
  },
}));

export default function ClustersPanel() {
  const classes = useStyles();

  const [zones, setZones] = React.useState(['white', 'blue', 'green'])
  const [requests, setRequests] = React.useState([
    { id: "id0", created: new Date().toDateString(), numberOfClusters: 100, requestedBy: "Some User", deleteInHours: 24, status: "provisioning"},
    { id: "id1", created: new Date().toDateString(), numberOfClusters: 100, requestedBy: "Some User", deleteInHours: 24, status: "provisioning"},
    { id: "id2", created: new Date().toDateString(), numberOfClusters: 100, requestedBy: "Some User", deleteInHours: 24, status: "provisioning"},
  ]);
  const [clusters, setClusters] = React.useState([
    { id: "id0", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id1", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id2", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id3", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id4", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id5", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id10", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id11", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id12", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id13", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id14", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id15", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id20", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id21", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id22", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id23", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id24", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
    { id: "id25", name: "name 0", url: "http://some.url/", status: "running", error: "Some Error Message"},
  ]);

  const onSelectRequest = (request) => {
    console.log("Selected Request: " + request.id);
  }

  const onSubmitRequest = (request) => {
    console.log(request);
  }

  const onDeleteCluster = (cluster) => {
    console.log(cluster);
  }

  return (
    <div style={{ display: 'flex', flex: 1, height: '100%' }}>
      <div className={classes.clusterpanel}>
        <div className={classes.form}>
            <RequestForm zones={zones} onSubmit={onSubmitRequest} />
        </div>
        <div className={classes.tables}>
          <div className={classes.table}>
            <RequestTable requests={requests} onSelect={onSelectRequest} />
          </div>
          <div className={classes.table}>
            <ClusterTable clusters={clusters} onDeleteCluster={onDeleteCluster} />
          </div>
        </div>
      </div>
    </div>
  );
}

