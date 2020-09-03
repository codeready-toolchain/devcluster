
// interval reference
var intervalRef;

var idToken;

// this is where we load our config from
configURL = '/api/v1/authconfig'

// loads json data from url, the callback is called with
// error and data, with data the parsed json.
var getJSON = function(method, url, token, params, callback) {
  var xhr = new XMLHttpRequest();
  xhr.open(method, url, true);
  if (token != null) {
    xhr.setRequestHeader('Authorization', 'Bearer ' + token)
  }
  xhr.responseType = 'json';
  xhr.onload = function() {
    var status = xhr.status;
    if (status >= 200 && status < 300) {
      callback(null, xhr.response);
    } else {
      callback(status, xhr.response);
    }
  };
  if (params != null) {
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.send(params)
  } else {
    xhr.send();
  }
};

// hides all state content.
function hideAll() {
  document.getElementById('state-waiting-for-provisioning').style.display = 'none';
  document.getElementById('state-request-clusters').style.display = 'none';
  document.getElementById('state-provisioned').style.display = 'none';
  document.getElementById('state-not-logged-in').style.display = 'none';
  document.getElementById('state-error').style.display = 'none';
  document.getElementById('dashboard').style.display = 'none';
}

// shows state content. Given Id needs to be one of the element
function show(elementId) {
  document.getElementById(elementId).style.display = 'block';
}

function showError(errorText) {
  hideAll();
  show('state-error');
  document.getElementById('errorStatus').textContent = errorText;
}

// shows a logged in user.
function showUser(username) {
  document.getElementById('username').textContent = username;
  document.getElementById('user-loggedin').style.display = 'inline';
  document.getElementById('user-notloggedin').style.display = 'none';
}

// shows login/signup button
function hideUser() {
  document.getElementById('username').textContent = '';
  document.getElementById('user-loggedin').style.display = 'none';
  document.getElementById('user-notloggedin').style.display = 'inline';
}

// this loads the js library at location 'url' dynamically and
// calls 'cbSuccess' when the library was loaded successfully
// and 'cbError' when there was an error loading the library.
function loadAuthLibrary(url, cbSuccess, cbError) {
  var script = document.createElement('script');
  script.setAttribute('src', url);
  script.setAttribute('type', 'text/javascript');
  var loaded = false;
  var loadFunction = function () {
    if (loaded) return;
    loaded = true;
    cbSuccess();
  };
  var errorFunction = function (error) {
    if (loaded) return;
    cbError(error)
  };
  script.onerror = errorFunction;
  script.onload = loadFunction;
  script.onreadystatechange = loadFunction;
  document.getElementsByTagName('head')[0].appendChild(script);
}
      
// gets the cluster requests once.
function getClusterRequests(cbSuccess, cbError) {
  getJSON('GET', '/api/v1/cluster-reqs', idToken, null,function(err, data) {
    if (err != null) {
      cbError(err, data);
    } else {
      cbSuccess(data);
    }
  })
}

// updates the provisioning state.
function updateProvisioningState() {
  getClusterRequests(function(data) {
      hideAll();
      show('dashboard')
      // Display all requests
  }, function(err, data) {
    if (err === 401) {
      // user is unauthorized, show login/signup view; stop interval.
      // clearInterval(intervalRef);
      hideUser();
      hideAll();
      show('state-not-logged-in');
      show('state-error');
      if(data != null && data.error != null){
        document.getElementById('errorStatus').textContent = data.error;
      }
    } else {
      // other error, show error box.
      showError(err);
    }
  })
}

function login() {
  keycloak.login()
}

// request cluster provisioning
function requestClusters() {
  var n = document.getElementById("number-of-clusters").value;
  getJSON('POST', '/api/v1/cluster-req', idToken, "number-of-clusters=" + n, function(err, data) {
    if (err != null) {
      showError(JSON.stringify(data, null, 2));
    } else {
      hideAll();
      show('state-waiting-for-provisioning');
    }
  });
  // intervalRef = setInterval(updateProvisioningState, 1000);
}
      
// main operation, load config, load client, run client
getJSON('GET', configURL, null, null, function(err, data) {
  if (err !== null) {
    console.log('error loading client config' + err);
    showError(err);
  } else {
    loadAuthLibrary(data['auth-client-library-url'], function() {
      console.log('client library load success!')
      var clientConfig = JSON.parse(data['auth-client-config']);
      console.log('using client configuration: ' + JSON.stringify(clientConfig))
      keycloak = Keycloak(clientConfig);
      keycloak.init({
        onLoad: 'check-sso',
        silentCheckSsoRedirectUri: window.location.origin + '/silent-check-sso.html',
      }).success(function(authenticated) {
        if (authenticated == true) {
          keycloak.loadUserInfo().success(function(data) {
            idToken = keycloak.idToken
            showUser(data.preferred_username)
            // now check the provisioning state
            //updateProvisioningState();
            hideAll();
            show('state-request-clusters');
          });
        } else {
          hideUser();
          hideAll();
          show('state-not-logged-in');
        }
      }).error(function() {
        console.log('Failed to initialize authorization');
        showError('Failed to initialize authorization.');
      });
    }, function(err) {
      console.log('error loading client library' + err);
      showError(err);
    });
  }
});
  