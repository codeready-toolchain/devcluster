import axios from 'axios';

var baseUrl = window.location.origin.startsWith('http://localhost')?'https://devcluster-alexeykazakov-stage.apps.sandbox-stage.gb17.p1.openshiftapps.com':window.location.origin;
//var baseUrl = 'https://devcluster-devcluster-dev.apps.sandbox-stage.gb17.p1.openshiftapps.com'

// gets zones
export const getZones = async () => {
  let resp = await axios({
    method: 'GET',
    url: baseUrl + '/api/v1/zones',
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}
  
// gets the cluster requests once.
export const getClusterRequests = async () => {
  let resp = await axios({
    method: 'GET',
    url: baseUrl + '/api/v1/cluster-reqs',
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}
  
// gets the cluster request.
export const getClusterRequest = async (id) => {
  var bodyFormData = new FormData();
  bodyFormData.append('id', id);
  let resp = await axios({
    method: 'GET',
    url: baseUrl + '/api/v1/cluster-req/' + id,
    data: bodyFormData,
    headers: {'Content-Type': 'multipart/form-data' },
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}

// deletes the cluster.
export const deleteCluster = async (id) => {
  var bodyFormData = new FormData();
  bodyFormData.append('id', id);
  let resp = await axios({
    method: 'DELETE',
    url: baseUrl + '/api/v1/cluster/' + id,
    data: bodyFormData,
    headers: {'Content-Type': 'multipart/form-data' },
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}

// requests clusters.
export const requestClusters = async (n, zone, deleteInHours) => {
  var bodyFormData = new FormData();
  bodyFormData.append('number-of-clusters', n);
  bodyFormData.append('zone', zone);
  bodyFormData.append('delete-in-hours', deleteInHours);
  let resp = await axios({
    method: 'POST',
    url: baseUrl + '/api/v1/cluster-req',
    data: bodyFormData,
    headers: {'Content-Type': 'multipart/form-data' },
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}

// requests users.
export const requestUsers = async (n, startIndex) => {
  var bodyFormData = new FormData();
  bodyFormData.append('number-of-users', n);
  bodyFormData.append('start-index', startIndex);
  let resp = await axios({
    method: 'POST',
    url: baseUrl + '/api/v1/users',
    data: bodyFormData,
    headers: {'Content-Type': 'multipart/form-data' },
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}
    
// get users
export const getUsers = async () => {
  let resp = await axios({
    method: 'GET',
    url: baseUrl + '/api/v1/users',
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}

// gets the all clusters in a zone
export const getClustersRequestsByZone = async (zoneID) => {
  var bodyFormData = new FormData();
  bodyFormData.append('zone', zoneID);
  let resp = await axios({
    method: 'GET',
    url: baseUrl + '/api/v1/clusters?zone=' + zoneID,
    data: bodyFormData,
    headers: {'Content-Type': 'multipart/form-data' },
  });
  if (resp.status >= 200 && resp.status < 300) {
    return Promise.resolve(resp.data);
  }
  else {
    return Promise.reject(new Error('' + resp.status + ' ' + resp.statusText));
  }
}
