import React from 'react';

import { ExportToCsv } from 'export-to-csv';

import { confirmAlert } from 'react-confirm-alert';
import 'react-confirm-alert/src/react-confirm-alert.css';

import { makeStyles } from '@material-ui/core/styles';
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';

import { getZones, getClusterRequests, getClusterRequest, deleteCluster, requestClusters } from './services/backend';

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
    padding: '10px',
  },
  tables: {
    display: 'flex',
    flex: 1,
    flexDirection: 'row',
    overflow: 'hidden',
  },
  table: {
    flex: 1,
    overflow: 'auto',
    padding: '10px',
  },
  modal: {
    position: 'absolute',
    left: '40%',
    top: '40%',
    width: '20%',
    height: '20%',
    backgroundColor: theme.palette.background.paper,
    padding: theme.spacing(2, 4, 3),
  },
}));

export default function ClustersPanel() {
  const classes = useStyles();

  const [zones, setZones] = React.useState([])
  const [requests, setRequests] = React.useState([]);
  const [selectedRequest, setSelectedRequest] = React.useState();
  const [clusters, setClusters] = React.useState([]);
  const [snackOpen, setSnackOpen] = React.useState(false);
  const [snackMessage, setSnackMessage] = React.useState();

  const handleSnackClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackOpen(false);
  };

  React.useEffect(() => {
    async function fetchData() {
      // fetch zones
      try {
        let zones = await getZones();
        setZones(zones);
      } catch (e) {
        console.error('error fetching zones', e.message);
        setSnackMessage('Error fetching zones: ' + e.message);
        setSnackOpen(true);
      }
      // fetch requests
      try {
        let requests = await getClusterRequests();
        setRequests(requests);
      } catch (e) {
        console.error('error fetching cluster requests', e.message);
        setSnackMessage('Error fetching cluster requests: ' + e.message);
        setSnackOpen(true);
      }  
    }
    fetchData();
  }, []);

  const onSelectRequest = async (request) => {
    try {
      setSelectedRequest(request);
      let requestWithClusters = await getClusterRequest(request.ID);
      setClusters(requestWithClusters.Clusters);
    } catch (e) {
      console.error('error fetching clusters', e.message);
      setSnackMessage('Error fetching clusters: ' + e.message);
      setSnackOpen(true);
    }
  }

  const onSubmitRequest = async (request) => {
    try {
      await requestClusters(request.numberOfClusters, request.zone, request.deleteInHours, request.subnet);
    } catch (e) {
      console.error('error requesting clusters', e.message);
      setSnackMessage('Error requesting clusters: ' + e.message);
      setSnackOpen(true);
    }
    // refresh requests
    try {
      let requests = await getClusterRequests();
      setRequests(requests);
    } catch (e) {
      console.error('error fetching cluster requests', e.message);
      setSnackMessage('Error fetching cluster requests: ' + e.message);
      setSnackOpen(true);
    }      
  }

  const onExportRequest = (request) => {      
      getClusterRequest(request.ID).then((result) => {
        let exportData = [];
        let clusters = result.Clusters;
        clusters.map((cluster) => {
          return exportData.push({
            'Cluster Id': cluster.ID,
            'Cluster Name': cluster.Name,
            'Master URL': cluster.MasterURL,
            'User Id': cluster.User.ID,
            'User E-Mail': cluster.User.Email,
            'User Password': cluster.User.Password,
            'User Policy Id': cluster.User.PolicyID,
            'User CloudDirect Id': cluster.User.CloudDirectID,
          });
        });
        const options = { 
          fieldSeparator: ',',
          quoteStrings: '"',
          decimalSeparator: '.',
          showLabels: true, 
          showTitle: false,
          useTextFile: false,
          useBom: true,
          useKeysAsHeaders: true,
        };
        const csvExporter = new ExportToCsv(options);
        csvExporter.generateCsv(exportData);
      })
  }

  const onDeleteCluster = (cluster) => {
    confirmAlert({
      title: 'Confirm to delete cluster',
      message: 'Cluster ' + cluster.ID + ' is being deleted. This can not be reverted. Are you sure.',
      buttons: [
        {
          label: 'Yes',
          onClick: () => onConfirmDeleteCluster(cluster),
        },
        {
          label: 'No',
        }
      ]
    });
  }

  const onConfirmDeleteCluster = async (cluster) => {
    try {
      await deleteCluster(cluster.ID);
      setSnackMessage('Cluster deleted..');
      setSnackOpen(true);
    } catch (e) {
      console.error('error deleting cluster', e.state, e.message);
      setSnackMessage('Error deleting cluster: ' + e.message);
      setSnackOpen(true);
    }
    try {
      let requestWithClusters = await getClusterRequest(selectedRequest.ID);
      setClusters(requestWithClusters.Clusters);
    } catch (e) {
      console.error('error fetching clusters', e.state, e.message);
      setSnackMessage('Error fetching clusters: ' + e.message);
      setSnackOpen(true);
    }
  }

  return (
    <div style={{ display: 'flex', flex: 1, height: '100%' }}>
      <div className={classes.clusterpanel}>
        <div className={classes.form}>
            <RequestForm zones={zones} onSubmit={onSubmitRequest} />
        </div>
        <div className={classes.tables}>
          <div className={classes.table}>
            <RequestTable requests={requests} onSelect={onSelectRequest} onExport={onExportRequest} />
          </div>
          <div className={classes.table}>
            <ClusterTable clusters={clusters} onDeleteCluster={onDeleteCluster} />
          </div>
        </div>
      </div>
      <Snackbar
        anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
        open={snackOpen}
        autoHideDuration={6000}
        onClose={handleSnackClose}
        message={snackMessage}
        action={
          <React.Fragment>
            <IconButton size='small' aria-label='close' color='inherit' onClick={handleSnackClose}>
              <CloseIcon fontSize='small' />
            </IconButton>
          </React.Fragment>
        }
      />
    </div>
  );
}

