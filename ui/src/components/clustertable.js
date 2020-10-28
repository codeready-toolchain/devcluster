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
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Collapse from '@material-ui/core/Collapse';
import FileCopyIcon from '@material-ui/icons/FileCopy';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import PasswordField from 'material-ui-password-field';
import { CopyToClipboard } from 'react-copy-to-clipboard';

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
    oneLineTable: {
      display: 'inline-block',
      whiteSpace: 'nowrap',
      overflow: 'hidden',
      textOverflow: 'ellipsis',
      maxWidth: 500,
    },
    copyFlex: {
      display: 'flex',
      alignItems: 'center'
    }
});

function Row(props) {
    const { row } = props;
    const { onSelect } = props;
    const { onDelete } = props;
    const { selected } = props;

    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();

    return (
      <React.Fragment>
        <TableRow className={classes.root} key={row.ID} hover onClick={(event) => onSelect(row)} aria-checked={selected} selected={selected}>
          <TableCell>
            <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
              {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
            </IconButton>
          </TableCell>
          <TableCell component="th" scope="row" className={classes.oneLine}>{row.ID.substring(row.ID.length - 5, row.ID.length)}</TableCell>
          <TableCell align="right">{row.Name}</TableCell>
          <TableCell align="right">{row.Status}</TableCell>
          <TableCell align="right">{row.Error?'ERROR':'No Error'}</TableCell>
          <TableCell>
            <IconButton aria-label='export' color='primary' onClick={() => onDelete(row)}>
              <DeleteIcon/>
            </IconButton>
          </TableCell>
        </TableRow>
        <TableRow>
          <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
            <Collapse in={open} timeout="auto" unmountOnExit>
              <Box margin={1}>
                <Typography variant="h6" gutterBottom component="div">Cluster Details</Typography>
                <Table>
                    <tbody>
                        <tr><td><Typography>Id:</Typography></td><td>{row.ID}</td></tr>
                        <tr><td><Typography>Name:</Typography></td><td>{row.Name}</td></tr>
                        <tr><td><Typography>Status:</Typography></td><td>{row.Status}</td></tr>
                        <tr><td><Typography>Error Message:</Typography></td><td>{row.Error?row.Error:'n/a'}</td></tr>
                        <tr>
                          <td><Typography>Hostname:</Typography></td>
                          <td className={classes.copyFlex}>
                            <div className={classes.oneLineTable}>{!row.Hostname?'n/a':row.Hostname}</div>
                            <CopyToClipboard text={row.Hostname}>
                              <IconButton className={classes.copyButton} size="small"><FileCopyIcon /></IconButton>
                            </CopyToClipboard>
                          </td>
                        </tr>
                        <tr>
                          <td><Typography>Console URL:</Typography></td>
                          <td className={classes.copyFlex}>
                            <div className={classes.oneLineTable}>{!row.ConsoleURL?'n/a':row.ConsoleURL}</div>
                            <CopyToClipboard text={row.ConsoleURL}>
                              <IconButton className={classes.copyButton} size="small"><FileCopyIcon /></IconButton>
                            </CopyToClipboard>
                          </td>
                        </tr>
                        <tr>
                          <td><Typography>Master URL:</Typography></td>
                          <td className={classes.copyFlex}>
                            <div className={classes.oneLineTable}>{!row.MasterURL?'n/a':row.MasterURL}</div>
                            <CopyToClipboard text={row.MasterURL}>
                              <IconButton className={classes.copyButton} size="small"><FileCopyIcon /></IconButton>
                            </CopyToClipboard>
                          </td>
                        </tr>
                        <tr>
                          <td><Typography>Login URL:</Typography></td>
                          <td className={classes.copyFlex}>
                            <div className={classes.oneLineTable}>{!row.LoginURL?'n/a':row.LoginURL}</div>
                            <CopyToClipboard text={row.LoginURL}>
                              <IconButton className={classes.copyButton} size="small"><FileCopyIcon /></IconButton>
                            </CopyToClipboard>
                          </td>
                        </tr>
                        <tr>
                          <td><Typography>Workshop URL:</Typography></td>
                          <td className={classes.copyFlex}>
                            <div className={classes.oneLineTable}>{!row.WorkshopURL?'n/a':row.WorkshopURL}</div>
                            <CopyToClipboard text={row.WorkshopURL}>
                              <IconButton className={classes.copyButton} size="small"><FileCopyIcon /></IconButton>
                            </CopyToClipboard>
                          </td>
                        </tr>
                        <tr><td><Typography>Username:</Typography></td><td>{row.User.ID}</td></tr>
                        <tr><td><Typography>User Password:</Typography></td><td><PasswordField visible={false} defaultValue={row.User.Password} inputProps={{readOnly: true,}}/></td></tr>
                    </tbody>
                </Table>
              </Box>
            </Collapse>
          </TableCell>
        </TableRow>
      </React.Fragment>
    );
  }
  
export default function ClusterTable({ clusters, onSelect, onDeleteCluster }) {

    const classes = useStyles();

    const [selectedCluster, setSelectedCluster] = React.useState();

    const handleClusterRowClick = (cluster) => {
        setSelectedCluster(cluster);
        if (onSelect)
            onSelect(cluster);
    }
    
    const isClusterSelected = (cluster) => 
        cluster && selectedCluster && cluster.ID === selectedCluster.ID;

    return (
        <TableContainer className={classes.container} component={Paper}>
        <Table stickyHeader className={classes.table} aria-label='cluster-table collapsible table'>
        <TableHead>
            <TableRow>
                <TableCell style={{width: '20px'}}/>
                <TableCell align="right" style={{width: '40px'}}>Id</TableCell>
                <TableCell align="right">Name</TableCell>
                <TableCell align="right">Status</TableCell>
                <TableCell align="right">Error</TableCell>
                <TableCell/>
            </TableRow>
        </TableHead>
        <TableBody>
            {clusters?clusters.map((cluster) => {
                return (<Row key={cluster.ID} row={cluster} selected={isClusterSelected(cluster)} onSelect={() => handleClusterRowClick(cluster)} onDelete={() => onDeleteCluster(cluster)}/>)
            }):null}
        </TableBody>
        </Table>
        </TableContainer>
    );
}
