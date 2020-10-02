import axios from 'axios';

function HttpException(status, message) {
    this.name = "HttpException";
    this.status = status;
    this.message = message;
}

// loads json data from url, the callback is called with
// error and data, with data the parsed json.
var getJSON = function(method, url, params, callback) {
  axios({
    method: method,
    url: url,
    data: params
  }).then(function (resp) {
    if (resp.status >= 200 && resp.status < 300) {
      callback(null, resp.data);
    } else {
      callback(status, resp.statusText);
    }
  });
}

// gets zones
export const getZones = async () => {
    getJSON('GET', '/api/v1/zones', null, function(err, data) {
      if (err != null) {
        throw new HttpException(err, JSON.stringify(data, null, 2));
      } else {
        return Promise.resolve(data);
      }
    })
}
  
// gets the cluster requests once.
export const getClusterRequests = async () => {
    getJSON('GET', '/api/v1/cluster-reqs', null,function(err, data) {
      if (err != null) {
        throw new HttpException(err, JSON.stringify(data, null, 2));
      } else {
        return Promise.resolve(data);
      }
    })
}
  
// gets the cluster request.
export const getClusterRequest = async (id) => {
    getJSON('GET', '/api/v1/cluster-req/' + id, null,function(err, data) {
      if (err != null) {
        throw new HttpException(err, JSON.stringify(data, null, 2));
      } else {
        Promise.resolve(data);
      }
    })
}

// deletes the cluster.
export const deleteCluster = async (id) => {
    getJSON('DELETE', '/api/v1/cluster/' + id, null,function(err, data) {
      if (err != null) {
        throw new HttpException(err, JSON.stringify(data, null, 2));
      } else {
        Promise.resolve(data);
      }
    })
}

// requests clusters.
export const requestClusters = async (n, zone, deleteInHours, noSubnet) => {
    getJSON('POST', '/api/v1/cluster-req', 
        "number-of-clusters=" + n + "&zone=" + zone + "&delete-in-hours=" + deleteInHours + noSubnet?"&no-subnet=true":"", 
        function(err, data) {
            if (err != null) {
                throw new HttpException(err, JSON.stringify(data, null, 2));
            } else {
                Promise.resolve(data);
            }
        })
}

// requests users.
export const requestUsers = async (n, startIndex) => {
    getJSON('POST', '/api/v1/users', 
        "number-of-users=" + n + "&start-index=" + startIndex, 
        function(err, data) {
            if (err != null) {
                throw new HttpException(err, JSON.stringify(data, null, 2));
            } else {
                Promise.resolve(data);
            }
        })
}

    
// get Users
export const getUsers = async () => {
    getJSON('GET', '/api/v1/users', "", function(err, data) {
      if (err != null) {
        throw new HttpException(err, JSON.stringify(data, null, 2));
      } else {
        Promise.resolve(data);
      }
    });
  }
  