import axios from 'axios';

var mockEnabled = false;

function HttpException(status, message) {
    this.name = "HttpException";
    this.status = status;
    this.message = message;
}

// loads json data from url, the callback is called with
// error and data, with data the parsed json.
var getJSON = function(method, url, params, callback) {
  console.log("request: " + method + " " + url  + " - " + JSON.stringify(params));
  axios({
    method: method,
    url: url,
    data: params
  }).then(function (resp) {
    if (resp.status >= 200 && resp.status < 300) {
      callback(null, resp.data);
    } else {
      callback(resp.status, resp.statusText);
    }
  });
}

// gets zones
export const getZones = async () => {
  if (mockEnabled) {
    console.log("request mock zones..");
    return Promise.resolve([{
      "id":           "lon06",
      "name":         "lon06",
      "kind":         "dc",
      "display_name": "London 06",
    },{
      "id":           "sng01",
      "name":         "sng01",
      "kind":         "dc",
      "display_name": "Singapore 01",
    }]);
  };
  getJSON('GET', '/api/v1/zones', null, function(err, data) {
    if (err != null) {
      throw new HttpException(err, JSON.stringify(data, null, 2));
    } else {
      return Promise.resolve(data);
    }
  });
}
  
// gets the cluster requests once.
export const getClusterRequests = async () => {
  if (mockEnabled) {
    console.log("request mock requests..");
    return Promise.resolve([{
      "ID":            "r0",
      "Requested":     42,
      "Created":       new Date().getTime(),
      "Status":        "ready",
      "Error":         "no error",
      "RequestedBy":   "Some User",
      "Zone":          "Singapore 01",
      "DeleteInHours": 72,
      "NoSubnet":      true,
    },{
      "ID":            "r1",
      "Requested":     99,
      "Created":       new Date().getTime(),
      "Status":        "ready",
      "Error":         "no error",
      "RequestedBy":   "Some User 2",
      "Zone":          "Singapore 02",
      "DeleteInHours": 48,
      "NoSubnet":      true,
    },]);
  };
  getJSON('GET', '/api/v1/cluster-reqs', null,function(err, data) {
    if (err != null) {
      throw new HttpException(err, JSON.stringify(data, null, 2));
    } else {
      return Promise.resolve(data);
    }
  });
}
  
// gets the cluster request.
export const getClusterRequest = async (id) => {
  if (mockEnabled) {
    console.log("request mock request details..");
    return Promise.resolve({
      "ID":            "r0",
      "Requested":     42,
      "Created":       new Date().getTime(),
      "Status":        "ready",
      "Error":         "no error",
      "RequestedBy":   "Some User",
      "Zone":          "Singapore 01",
      "DeleteInHours": 72,
      "NoSubnet":      true,
      "Clusters":      [{
        "ID":        "c0",
        "RequestID": "r1",
        "Name":      "Cluster Name 0",
        "URL":       "http://cluster0.url/",
        "Status":    "ready",
        "Error":     "no error",
        "User":      {
          "ID":            "u0",
          "CloudDirectID": "cdid0",
          "Email":         "some0@user.com",
          "Password":      "secret",
          "ClusterID":     "c0",
          "PolicyID":      "PolicyId",
          "Recycled":      new Date().getTime(),
        }
      },{
        "ID":        "c1",
        "RequestID": "r1",
        "Name":      "Cluster Name 1",
        "URL":       "http://cluster1.url/",
        "Status":    "ready",
        "Error":     "no error",
        "User":      {
          "ID":            "u1",
          "CloudDirectID": "cdid0",
          "Email":         "some1@user.com",
          "Password":      "secret",
          "ClusterID":     "c1",
          "PolicyID":      "PolicyId",
          "Recycled":      new Date().getTime(),
        }
      }]
    });
  };
  getJSON('GET', '/api/v1/cluster-req/' + id, null,function(err, data) {
    if (err != null) {
      throw new HttpException(err, JSON.stringify(data, null, 2));
    } else {
      Promise.resolve(data);
    }
  });
}

// deletes the cluster.
export const deleteCluster = async (id) => {
  if (mockEnabled) {
    console.log("request cluster delete, id=" + id);
    return Promise.resolve();
  }
  getJSON('DELETE', '/api/v1/cluster/' + id, null,function(err, data) {
    if (err != null) {
      throw new HttpException(err, JSON.stringify(data, null, 2));
    } else {
      Promise.resolve(data);
    }
  });
}

// requests clusters.
export const requestClusters = async (n, zone, deleteInHours, noSubnet) => {
  if (mockEnabled) {
    console.log("request cluster, n=" + n + " zone=" + zone + " deleteInHours=" + deleteInHours + " noSubnet=" + noSubnet);
    return Promise.resolve();
  }
  getJSON('POST', '/api/v1/cluster-req', 
      "number-of-clusters=" + n + "&zone=" + zone + "&delete-in-hours=" + deleteInHours + noSubnet?"&no-subnet=true":"", 
      function(err, data) {
          if (err != null) {
              throw new HttpException(err, JSON.stringify(data, null, 2));
          } else {
              Promise.resolve(data);
          }
  });
}

// requests users.
export const requestUsers = async (n, startIndex) => {
  if (mockEnabled) {
    console.log("request users, n=" + n + " startIndex=" + startIndex);
    return Promise.resolve();
  }
  getJSON('POST', '/api/v1/users', 
      "number-of-users=" + n + "&start-index=" + startIndex, 
      function(err, data) {
          if (err != null) {
              throw new HttpException(err, JSON.stringify(data, null, 2));
          } else {
              Promise.resolve(data);
          }
  });
}
    
// get Users
export const getUsers = async () => {
  if (mockEnabled) {
    console.log("request users..");
    return Promise.resolve([{
      "ID":            "u0",
      "CloudDirectID": "cdid0",
      "Email":         "some0@user.com",
      "Password":      "secret",
      "ClusterID":     "c0",
      "PolicyID":      "PolicyId",
      "Recycled":      new Date().getTime(),
    }, {
      "ID":            "u1",
      "CloudDirectID": "cdid0",
      "Email":         "some1@user.com",
      "Password":      "secret",
      "ClusterID":     "c0",
      "PolicyID":      "PolicyId",
      "Recycled":      new Date().getTime(),
    }]);
  }
  getJSON('GET', '/api/v1/users', "", function(err, data) {
    if (err != null) {
      throw new HttpException(err, JSON.stringify(data, null, 2));
    } else {
      Promise.resolve(data);
    }
  });
}
  