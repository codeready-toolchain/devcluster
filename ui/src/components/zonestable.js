import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import { LinearProgress } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
    container: {
        height: '100%',
    },
    oneLine: {
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis',  
    },
    barOn: {
      display: 'block',
    },
    barOff: {
      display: 'none',
    },
}));

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
    const { selected } = props;
    const { onSelect } = props;

    const classes = useRowStyles();

    let rowDate = new Date(0);
    rowDate.setUTCSeconds(row.Recycled);

    return (
      <React.Fragment>
        <TableRow className={classes.root} key={row.ID} hover onClick={(event) => onSelect(row)} aria-checked={selected} selected={selected}>
          <TableCell align="right">{row.zoneID}</TableCell>
          <TableCell align="right">{row.zoneName}</TableCell>
          <TableCell align="right">{row.activeClusters}</TableCell>
        </TableRow>
      </React.Fragment>
    );
  }
  
export default function ZonesTable({ zones, inProgress, onSelect }) {

    const classes = useStyles();

    const [selectedZone, setSelectedZone] = React.useState();

    const handleZoneRowClick = (zone) => {
        setSelectedZone(zone);
        if (onSelect)
            onSelect(zone);
    }
    
    const isZoneSelected = (zone) => 
      zone && selectedZone && zone.zoneID === selectedZone.zoneID;

    return (
        <TableContainer className={classes.container} component={Paper}>
          <LinearProgress classes={{ bar: inProgress?classes.barOn:classes.barOff}} />
          <Table stickyHeader className={classes.table} aria-label='zone-table collapsible table'>
            <TableHead>
                <TableRow>
                    <TableCell align="right" style={{width: '40px'}}>Id</TableCell>
                    <TableCell align="right">Name</TableCell>
                    <TableCell align="right">Active Clusters Total</TableCell>
                    <TableCell/>
                </TableRow>
            </TableHead>
            <TableBody>
                {zones?zones.map((zone) => {
                    return (<Row key={zone.zoneID} row={zone} selected={isZoneSelected(zone)} onSelect={() => handleZoneRowClick(zone)}/>)
                }):null}
            </TableBody>
          </Table>
        </TableContainer>
    );
}
