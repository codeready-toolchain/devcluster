import React from 'react';

import axios from 'axios';

import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import ClustersPanel from './clusterspanel';
import UsersPanel from './userspanel';

import logo from './redhat-logo.svg';
import rhdlogo from './rhdeveloper-logo.svg';

const useStyles = makeStyles((theme) => ({
  app: {
    display: 'flex',
    flex: 1,
    flexDirection: "column",
    height: "100%",
  },
  appbar: {
    background: "black",
  },
  tabPanel: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: "column",
    overflow: 'auto',
    padding: "15px",
  },
  tabPanelCentered: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: "column",
    overflow: 'auto',
    padding: "15px",
    alignItems: 'center',
    justifyContent: 'center',
  },
  title: {
    flexGrow: 1,
    paddingLeft: "25px",
    paddingTop: "3px",
  },
  toplogo: {
    height: "35px",
  },
  userContainer: {
    display: 'flex',
    flexDirection: 'row',
    flexGrow: 1,
    alignItems: 'center',
    justifyContent: 'flex-end',

  },
  username: {
    paddingRight: '10px',
  },
  loginRequiredPanel: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  centerLogo: {
    height: '100px',
    marginBottom: '20px'
  },
  centerButton: {
    marginTop: '10px'
  },
}));

export default function App() {
  const classes = useStyles();

  const [activeTab, setActiveTab] = React.useState("tab-clusters");

  const handleChange = (event, newValue) => {
    setActiveTab(newValue);
  };

  const [keycloak, setKeycloak] = React.useState();
  const [authenticated, setAuthenticated] = React.useState(false);
  const [username, setUsername] = React.useState();

  const attachClientElem = (url, callback) => {
    console.log('outer layer, using keycloak client url: ' + url);
    if (!document.getElementById('kc-client-script')) {
      const script = document.createElement('script');
      script.type = "text/javascript"
      script.src = url; 
      script.id = 'kc-client-script'; 
      document.head.appendChild(script);
      script.onload = () => {
        console.log('outer layer loaded keycloak client script');
        if (callback)
          callback();
      };
    }
  }

  const loadKeycloakClient = (callback) => {
    let url;
    if (window.location.origin.startsWith('http://localhost')) {
      url = 'https://sso.prod-preview.openshift.io/auth/js/keycloak.js';
      attachClientElem(url, callback);
    } else {
      fetch(window.location.origin + '/api/v1/authconfig')
        .then(response => response.json())
        .then((jsonData) => {
          console.log('outer layer, fetched auth config: ' + JSON.stringify(jsonData));
          if (!jsonData['auth-client-library-url']) {
            throw new Error('outer layer, loaded auth config is malformed!'); 
          }
          url = jsonData['auth-client-library-url'];
          attachClientElem(url, callback);
        }).catch((error) => {
          console.error('outer layer, error fetching keycloak config: ' + error)
        })
    }
  }

  const initKeycloak = (keycloakClient) => {
    keycloakClient.init({onLoad: 'check-sso', silentCheckSsoRedirectUri: window.location.origin})
    .success(authenticated => {
      axios.defaults.headers.common['Authorization'] = 'Bearer ' + keycloakClient.idToken;
      setKeycloak(keycloakClient);
      setAuthenticated(authenticated);
      keycloakClient.loadUserInfo().success(function(data) {
        setUsername(data.preferred_username)
      });
      console.log('keycloak init complete');
    })
    .error((error) => {
      console.warn('keycloak client init failed:', error);
    });
  }

  React.useEffect(() => {
    loadKeycloakClient(() => {
      console.log('client loaded, starting auth init..');
      const Keycloak = window.Keycloak;
      var keycloakClient;
      if (window.location.origin.startsWith('http://localhost')) {
        keycloakClient = new Keycloak("./keycloak.json");
        console.log('keycloak create complete, using local config');
        initKeycloak(keycloakClient);
      } else {
        fetch(window.location.origin + '/api/v1/authconfig')
          .then(response => response.json())
          .then((jsonData) => {
            console.log('fetched auth config: ' + JSON.stringify(jsonData));
            if (!jsonData["auth-client-config"]) {
              throw new Error('loaded keycloak config is malformed!'); 
            }
            console.log('using keycloak config: ' + JSON.stringify(jsonData["auth-client-config"]));
            keycloakClient = new Keycloak(JSON.parse(jsonData["auth-client-config"]));
            console.log('keycloak create complete, using remote config');
            initKeycloak(keycloakClient);
        }).catch((error) => {
          console.error('error fetching keycloak config: ' + error)
        });
      }
    });
  }, []);
  
  return (
      <div className={classes.app}>
        <AppBar position="static" className={classes.appbar}>
          <Toolbar>
            <img src={logo} className={classes.toplogo} alt="Red Hat" />
            <Typography variant="h6" align="left" className={classes.title}>Dev Clusters Dashboard</Typography>
            { !authenticated && <Button variant="contained" color="primary" onClick={() => keycloak.login()}>Login</Button> }
            { authenticated && 
              <div className={classes.userContainer}>
                <Typography align="left" className={classes.username}>{username}</Typography>
                <Button variant="contained" color="primary" onClick={() => keycloak.logout()}>Logout</Button>
              </div>
            }
          </Toolbar>
        </AppBar>
        { !authenticated && 
          <div className={classes.loginRequiredPanel}>
            <img src={rhdlogo} className={classes.centerLogo} alt="Red Hat" />
            <Typography>Please log in to the Dev Clusters Dashboard</Typography>
            <Button variant="contained" color="primary" className={classes.centerButton} onClick={() => keycloak.login()}>Login</Button>
          </div>
        }
        { authenticated &&
          <AppBar position="static" color="default">
            <Tabs aria-label="main-tabs" value={activeTab} onChange={handleChange} >
              <Tab label="Clusters" value="tab-clusters" />
              <Tab label="Users" value="tab-users" />
            </Tabs>
          </AppBar>
        } 
        {authenticated && activeTab === 'tab-clusters' ? <div className={classes.tabPanel}><ClustersPanel key="tab-clusters" /></div> : null}
        {authenticated && activeTab === 'tab-users' ? <div className={classes.tabPanelCentered}><UsersPanel key="tab-users" /></div> : null}
      </div>
  );
}
