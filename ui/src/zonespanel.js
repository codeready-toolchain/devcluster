import React from 'react';

import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';

import { getZones, getClustersRequestsByZone } from './services/backend';

import ZonesTable from './components/zonestable';

export default function ZonesPanel() {

  const [zonesDetails, setZonesDetails] = React.useState([]);
  const [inProgress, setInProgress] = React.useState(false);
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
      setInProgress(true);
      try {
        let zones = await getZones();
        let zonesDetails = [];
        for (let i=0; i<zones.length; i++) {
          let clusters = await getClustersRequestsByZone(zones[i].id);
          zonesDetails.push({
            zoneID: zones[i]['id'],
            zoneName: zones[i]['display_name'],
            activeClusters: clusters.length,
          })
        }
        setInProgress(false);
        setZonesDetails(zonesDetails);
      } catch (e) {
        console.error('error fetching zones data', e.message);
        setSnackMessage('Error fetching zones data: ' + e.message);
        setSnackOpen(true);
      }
    }
    fetchData();
  }, []);

  const onSelectZone = async (zone) => {
  }

  return (
    <div style={{ display: 'flex', flex: 1, height: '100%', width: '100%' }}>
      <ZonesTable zones={zonesDetails} inProgress={inProgress} onSelect={onSelectZone} />
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

