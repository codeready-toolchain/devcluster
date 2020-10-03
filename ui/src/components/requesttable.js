import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import IconButton from '@material-ui/core/IconButton';
import CloudDownloadIcon from '@material-ui/icons/CloudDownload';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

const useStyles = makeStyles((theme) => ({
}));
  
export default function RequestTable({ requests, onSelect }) {

    const classes = useStyles();

    const [selectedRequest, setSelectedRequest] = React.useState();

    const handleRequestRowClick = (event, request) => {
        setSelectedRequest(request);
        onSelect(request);
    }
    
    const isRequestSelected = (request) => 
        request && selectedRequest && request.ID === selectedRequest.ID;

    return (
        <TableContainer component={Paper}>
        <Table className={classes.table} aria-label="request-table">
        <TableHead>
            <TableRow>
            <TableCell>Id</TableCell>
            <TableCell>Created</TableCell>
            <TableCell>#&nbsp;of&nbsp;Clusters</TableCell>
            <TableCell>Requested&nbsp;by</TableCell>
            <TableCell>Delete&nbsp;in&nbsp;hours</TableCell>
            <TableCell>Status</TableCell>
            <TableCell></TableCell>
            </TableRow>
        </TableHead>
        <TableBody>
            {requests.map((request) => (
            <TableRow 
                key={request.ID} 
                hover 
                onClick={(event) => handleRequestRowClick(event, request)}
                aria-checked={isRequestSelected(request)}
                selected={isRequestSelected(request)}>
                <TableCell component="th" scope="row">{request.ID}</TableCell>
                <TableCell>{new Date(request.Created).toDateString()}</TableCell>
                <TableCell align="right">{request.Requested}</TableCell>
                <TableCell>{request.RequestedBy}</TableCell>
                <TableCell align="right">{request.DeleteInHours}</TableCell>
                <TableCell>{request.Status}</TableCell>
                <TableCell>
                    <IconButton aria-label="export" color="primary">
                        <CloudDownloadIcon/>
                    </IconButton>
                </TableCell>
            </TableRow>
            ))}
        </TableBody>
        </Table>
        </TableContainer>
    );
}

