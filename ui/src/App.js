import React from 'react';

import axios from 'axios';

import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import ClustersPanel from './clusterspanel'

import logo from './redhat-logo.svg';

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
  title: {
    flexGrow: 1,
    paddingLeft: "25px",
    paddingTop: "3px",
  },
  toplogo: {
    height: "35px",
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
  
  React.useEffect(() => {
    const Keycloak = window.Keycloak;
    const keycloakClient = new Keycloak("./keycloak.json");
    keycloakClient.init({onLoad: 'check-sso', silentCheckSsoRedirectUri: window.location.origin})
      .success(authenticated => {
        setKeycloak(keycloakClient);
        setAuthenticated(authenticated);
        axios.defaults.headers.common['Authorization'] = 'Bearer ' + keycloakClient.idToken;
        keycloakClient.loadUserInfo().success(function(data) {
          setUsername(data.preferred_username)
        });
      }) 
      .error((error) => {
        console.warn('Keycloak client init failed:', error);
      });
  }, []);
  
  return (
      <div className={classes.app}>
        <AppBar position="static" className={classes.appbar}>
          <Toolbar>
            <img src={logo} className={classes.toplogo} alt="Red Hat" />
            <Typography variant="h6" align="left" className={classes.title}>Dev Clusters Dashboard</Typography>
            { !authenticated && <Button color="inherit" onClick={() => keycloak.login()}>Login</Button> }
            { authenticated && 
              <div>
                <Typography align="left" className={classes.username}>{username} <Button color="inherit" onClick={() => keycloak.logout()}>Logout</Button></Typography>
              </div>
            }
          </Toolbar>
        </AppBar>
        <AppBar position="static" color="default">
          <Tabs aria-label="main-tabs" value={activeTab} onChange={handleChange} >
            <Tab label="Clusters" value="tab-clusters" />
            <Tab label="Users" value="tab-users" />
          </Tabs>
        </AppBar>
        {activeTab === 'tab-clusters' ? <div className={classes.tabPanel}><ClustersPanel key="tab-clusters" /></div> : null}
        {activeTab === 'tab-users' ? <div className={classes.tabPanel}><ClustersPanel key="tab-users" /></div> : null}
      </div>
  );
}
