import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton';
import DeleteIcon from '@material-ui/icons/Delete';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

const useStyles = makeStyles((theme) => ({
    container: {
        height: "100%",
    },
}));
  
export default function ClusterTable({ clusters, onSelect, onDeleteCluster }) {

    const classes = useStyles();

    const [selectedCluster, setSelectedCluster] = React.useState();

    const handleClusterRowClick = (event, cluster) => {
        setSelectedCluster(cluster);
        if (onSelect)
            onSelect(cluster);
    }
    
    const isClusterSelected = (cluster) => 
        cluster && selectedCluster && cluster.ID === selectedCluster.ID;

    return (
        <TableContainer className={classes.container} component={Paper}>
        <Table stickyHeader className={classes.table} aria-label="cluster-table">
        <TableHead>
            <TableRow>
            <TableCell>Id</TableCell>
            <TableCell>Name</TableCell>
            <TableCell>URL</TableCell>
            <TableCell>Status</TableCell>
            <TableCell>Error</TableCell>
            <TableCell></TableCell>
            </TableRow>
        </TableHead>
        <TableBody>
            {clusters.map((cluster) => (
            <TableRow 
                key={cluster.ID} 
                hover 
                onClick={(event) => handleClusterRowClick(event, cluster)}
                aria-checked={isClusterSelected(cluster)}
                selected={isClusterSelected(cluster)}>
                <TableCell component="th" scope="row">{cluster.ID}</TableCell>
                <TableCell >{cluster.Name}</TableCell>
                <TableCell align="right">{cluster.URL}</TableCell>
                <TableCell >{cluster.Status}</TableCell>
                <TableCell align="right">{cluster.Error}</TableCell>
                <TableCell >
                    <IconButton aria-label="export" color="primary" onClick={() => onDeleteCluster(cluster)}>
                        <DeleteIcon/>
                    </IconButton>
                </TableCell>
            </TableRow>
            ))}
        </TableBody>
        </Table>
        </TableContainer>
    );
}

