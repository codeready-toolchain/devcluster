import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton';
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

const useStyles = makeStyles((theme) => ({
    container: {
        height: '100%',
    },
    oneLine: {
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis',  
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

    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

    let rowDate = new Date(0);
    rowDate.setUTCSeconds(row.Recycled);

    return (
      <React.Fragment>
        <TableRow className={classes.root} key={row.ID} hover onClick={(event) => onSelect(row)} aria-checked={selected} selected={selected}>
          <TableCell>
            <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
              {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
            </IconButton>
          </TableCell>
          <TableCell component="th" scope="row" className={classes.oneLine}>{row.ID}</TableCell>
          <TableCell align="right">{row.Email}</TableCell>
          <TableCell align="right">{row.ClusterID}</TableCell>
        </TableRow>
        <TableRow>
          <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={4}>
            <Collapse in={open} timeout="auto" unmountOnExit>
              <Box margin={1}>
                <Typography variant="h6" gutterBottom component="div">User Details</Typography>
                <Table>
                    <tbody>
                        <tr><td><Typography>Id:</Typography></td><td>{row.ID}</td></tr>
                        <tr><td><Typography>CloudDirect Id:</Typography></td><td>{row.CloudDirectID}</td></tr>
                        <tr><td><Typography>E-Mail:</Typography></td><td>{row.Email}</td></tr>
                        <tr><td><Typography>Password:</Typography></td><td>{row.Password}</td></tr>
                        <tr><td><Typography>Policy Id:</Typography></td><td>{row.PolicyID}</td></tr>
                        <tr><td><Typography>Cluster Id:</Typography></td><td>{row.ClusterID}</td></tr>
                        <tr><td><Typography>Last Recycled:</Typography></td><td>{rowDate.toString()}</td></tr>
                    </tbody>
                </Table>
              </Box>
            </Collapse>
          </TableCell>
        </TableRow>
      </React.Fragment>
    );
  }
  
export default function UserTable({ users, onSelect }) {

    const classes = useStyles();

    const [selectedUser, setSelectedUser] = React.useState();

    const handleUserRowClick = (user) => {
        setSelectedUser(user);
        if (onSelect)
            onSelect(user);
    }
    
    const isUserSelected = (user) => 
    user && selectedUser && user.ID === selectedUser.ID;

    return (
        <TableContainer className={classes.container} component={Paper}>
        <Table stickyHeader className={classes.table} aria-label='user-table collapsible table'>
        <TableHead>
            <TableRow>
                <TableCell style={{width: '20px'}}/>
                <TableCell align="right" style={{width: '40px'}}>Id</TableCell>
                <TableCell align="right">E-Mail</TableCell>
                <TableCell align="right">Cluster Id</TableCell>
                <TableCell/>
            </TableRow>
        </TableHead>
        <TableBody>
            {users?users.map((user) => {
                return (<Row key={user.ID} row={user} selected={isUserSelected(user)} onSelect={() => handleUserRowClick(user)}/>)
            }):null}
        </TableBody>
        </Table>
        </TableContainer>
    );
}
