import React from 'react';

import { ExportToCsv } from 'export-to-csv';

import { makeStyles } from '@material-ui/core/styles';
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';

import { getUsers, requestUsers } from './services/backend';

import UserForm from './components/userform';
import UserTable from './components/usertable';

const useStyles = makeStyles((theme) => ({
  userspanel: {
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

export default function UsersPanel() {
  const classes = useStyles();

  const [users, setUsers] = React.useState([]);
  const [selectedUser, setSelectedUser] = React.useState();
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
      // fetch users
      try {
        let users = await getUsers();
        setUsers(users);
      } catch (e) {
        console.error('error fetching user requests', e.message);
        setSnackMessage('Error fetching user requests: ' + e.message);
        setSnackOpen(true);
      }  
    }
    fetchData();
  }, []);

  const onSelectUser = async (user) => {
      setSelectedUser(user);
  }

  const onSubmitRequest = async (request) => {
    try {
      await requestUsers(request.numberOfUsers, request.startIndex);
    } catch (e) {
      console.error('error requesting users', e.message);
      setSnackMessage('Error requesting users: ' + e.message);
      setSnackOpen(true);
    }
    // refresh requests
    try {
      let users = await getUsers();
      setUsers(users);
    } catch (e) {
      console.error('error fetching user requests', e.message);
      setSnackMessage('Error fetching user requests: ' + e.message);
      setSnackOpen(true);
    }      
  }

  const onExportUsers = () => {      
      let exportData = [];
      users.map((user) => {
        exportData.push({
          'User Id': user.ID,
          'User E-Mail': user.Email,
          'User Password': user.Password,
          'User Policy Id': user.PolicyID,
          'User CloudDirect Id': user.CloudDirectID,
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
  }

  return (
    <div style={{ display: 'flex', flex: 1, height: '100%', width: '100%' }}>
      <div className={classes.userspanel}>
        <div className={classes.form}>
            <UserForm onSubmit={onSubmitRequest} onExport={onExportUsers}/>
        </div>
        <div className={classes.tables}>
          <div className={classes.table}>
            <UserTable users={users} onSelect={onSelectUser} />
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

