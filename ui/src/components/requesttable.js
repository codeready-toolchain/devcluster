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
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Collapse from '@material-ui/core/Collapse';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
  
const useRowStyles = makeStyles({
    root: {
      '& > *': {
        borderBottom: 'unset',
      },
    },
    oneLine: {
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis',  
    },
});

function Row(props) {
    const { row } = props;
    const { onSelect } = props;
    const { onExport } = props;
    const { selected } = props;

    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();
  
    let rowDate = new Date(0);
    rowDate.setUTCSeconds(row.Created);
    let rowDateStr = (rowDate.getMonth()+1) + "-" + rowDate.getDate() + "-" + rowDate.getFullYear();

    return (
      <React.Fragment>
        <TableRow className={classes.root} key={row.ID} hover onClick={() => onSelect(row)} aria-checked={selected} selected={selected}>
          <TableCell>
            <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
              {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
            </IconButton>
          </TableCell>
          <TableCell component="th" scope="row" className={classes.oneLine}>{row.ID.substring(row.ID.length - 5, row.ID.length)}</TableCell>
          <TableCell align="right">{rowDateStr}</TableCell>
          <TableCell align="right">{row.Requested}</TableCell>
          <TableCell align="right">{row.RequestedBy}</TableCell>
          <TableCell align="right">{row.DeleteInHours}</TableCell>
          <TableCell align="right">{row.Status}</TableCell>
          <TableCell>
            <IconButton aria-label='export' color='primary' onClick={() => onExport(row)}>
              <CloudDownloadIcon/>
            </IconButton>
          </TableCell>
        </TableRow>
        <TableRow>
          <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={8}>
            <Collapse in={open} timeout="auto" unmountOnExit>
              <Box margin={1}>
                <Typography variant="h6" gutterBottom component="div">Request Details</Typography>
                <Table>
                    <tbody>
                        <tr><td><Typography>Request Id:</Typography></td><td>{row.ID}</td></tr>
                        <tr><td><Typography>Number of Clusters:</Typography></td><td>{row.Requested}</td></tr>
                        <tr><td><Typography>No Subnet:</Typography></td><td>{row.NoSubnet?'true':'false'}</td></tr>
                        <tr><td><Typography>Zone:</Typography></td><td>{row.Zone}</td></tr>
                        <tr><td><Typography>Requested at:</Typography></td><td>{rowDate.toUTCString()}</td></tr>
                        <tr><td><Typography>Requested by:</Typography></td><td>{row.RequestedBy}</td></tr>
                        <tr><td><Typography>Deletes in hours:</Typography></td><td>{row.DeleteInHours}</td></tr>
                        <tr><td><Typography>Status:</Typography></td><td>{row.Status}</td></tr>
                    </tbody>
                </Table>
              </Box>
            </Collapse>
          </TableCell>
        </TableRow>
      </React.Fragment>
    );
  }
  

export default function RequestTable({ requests, onSelect, onExport }) {

    const [selectedRequest, setSelectedRequest] = React.useState();

    const handleRequestRowClick = (request) => {
        setSelectedRequest(request);
        onSelect(request);
    }
    
    const isRequestSelected = (request) => 
        request && selectedRequest && request.ID === selectedRequest.ID;

    return (
        <TableContainer component={Paper} style={{ overflow: "auto" }}>
            <Table style={{ tableLayout: "fixed" }} aria-label='request-table collapsible table'>
                <TableHead>
                    <TableRow>
                        <TableCell style={{width: '20'}}/>
                        <TableCell align="right" style={{width: '40px'}}>Id</TableCell>
                        <TableCell align="right">Created</TableCell>
                        <TableCell align="right">#&nbsp;of&nbsp;Clusters</TableCell>
                        <TableCell align="right">Requested&nbsp;by</TableCell>
                        <TableCell align="right">Delete&nbsp;in&nbsp;hours</TableCell>
                        <TableCell align="right">Status</TableCell>
                        <TableCell/>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {requests?requests.map((request) => {
                        return (<Row key={request.ID} row={request} selected={isRequestSelected(request)} onSelect={() => handleRequestRowClick(request)} onExport={() => onExport(request)}/>)
                    }):null}
                </TableBody>
            </Table>
        </TableContainer>
    );
}
