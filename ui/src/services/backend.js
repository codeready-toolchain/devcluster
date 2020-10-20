import axios from 'axios';

var baseUrl = 'https://devcluster-alexeykazakov-stage.apps.member.crt-stage.com';

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
export const requestClusters = async (n, zone, deleteInHours, noSubnet) => {
  var bodyFormData = new FormData();
  bodyFormData.append('number-of-clusters', n);
  bodyFormData.append('zone', zone);
  bodyFormData.append('delete-in-hours', deleteInHours);
  bodyFormData.append('no-subnet', (noSubnet?true:false));
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
  