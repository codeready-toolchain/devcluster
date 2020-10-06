import axios from 'axios';

var baseUrl = "https://devcluster-alexeykazakov-stage.apps.member.crt-stage.com";

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
  let resp = await axios({
    method: 'GET',
    url: baseUrl + '/api/v1/cluster-req',
    params: { 'id': id },
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
  let resp = await axios({
    method: 'DELETE',
    url: baseUrl + '/api/v1/cluster',
    params: { 'id': id },
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
  let resp = await axios({
    method: 'POST',
    url: baseUrl + '/api/v1/cluster-req',
    data: {
      "number-of-clusters": '' + n,
      "zone": zone,
      "delete-in-hours": deleteInHours,
      "no-subnet": (noSubnet?true:false)
    },
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
  let resp = await axios({
    method: 'POST',
    url: baseUrl + '/api/v1/cluster-req',
    params: {
      "number-of-users": n,
      "start-index": startIndex,
    }
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
  