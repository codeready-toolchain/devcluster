import React from 'react';
import clsx from 'clsx';

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
import TableSortLabel from '@material-ui/core/TableSortLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Toolbar from '@material-ui/core/Toolbar';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Collapse from '@material-ui/core/Collapse';
import FileCopyIcon from '@material-ui/icons/FileCopy';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import PasswordField from 'material-ui-password-field';
import { CopyToClipboard } from 'react-copy-to-clipboard';

const useStyles = makeStyles({
    container: {
        height: '100%',
    },
    oneLine: {
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis',  
    },
    visuallyHidden: {
      border: 0,
      clip: 'rect(0 0 0 0)',
      height: 1,
      margin: -1,
      overflow: 'hidden',
      padding: 0,
      position: 'absolute',
      top: 20,
      width: 1,
    },
});

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

const useToolbarStyles = makeStyles({
  highlight: {
    color: 'red',
    backgroundColor: 'lightred',
  },
  title: {
    flex: '1 1 100%',
  },
});

const headCells = [
  { id: 'ID', numeric: false, disablePadding: false, label: 'Id' },
  { id: 'Name', numeric: false, disablePadding: false, label: 'Name' },
  { id: 'Status', numeric: false, disablePadding: false, label: 'Status' },
  { id: 'Error', numeric: false, disablePadding: false, label: 'Error' },
];

function Row(props) {
  const { row, onSelect, selected } = props;
  const [open, setOpen] = React.useState(false);
  const classes = useRowStyles();
  return (
    <React.Fragment>
      <TableRow className={classes.root} key={row.ID} hover onClick={() => onSelect(row)} selected={selected}>
        <TableCell padding="checkbox">
          <Checkbox
            checked={selected}
            onClick={() => onSelect(row)}
          />
        </TableCell>
        <TableCell>
          <IconButton size="small" onClick={(event) => {event.stopPropagation(); setOpen(!open);}}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row" className={classes.oneLine}>{row.ID.substring(row.ID.length - 5, row.ID.length)}</TableCell>
        <TableCell align="left">{row.Name}</TableCell>
        <TableCell align="left">{row.Status}</TableCell>
        <TableCell align="left">{row.Error?'ERROR':'No Error'}</TableCell>
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
                      <tr>
                          <td><Typography>Username:</Typography></td>
                          <td className={classes.copyFlex}>
                              {row.User.ID}
                              <CopyToClipboard text={row.User.ID}>
                                  <IconButton className={classes.copyButton} size="small"><FileCopyIcon /></IconButton>
                              </CopyToClipboard>
                          </td>
                      </tr>
                      <tr>
                          <td><Typography>User Password:</Typography></td>
                          <td className={classes.copyFlex}>
                              <PasswordField visible={false} defaultValue={row.User.Password} inputProps={{readOnly: true,}}/>
                              <CopyToClipboard text={row.User.Password}>
                                  <IconButton className={classes.copyButton} size="small"><FileCopyIcon /></IconButton>
                              </CopyToClipboard>
                          </td>
                      </tr>
                  </tbody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
}

function EnhancedTableHead(props) {
  const { classes, onSelectAllClick, order, orderBy, numSelected, rowCount, onRequestSort } = props;
  const createSortHandler = (property) => (event) => {
    onRequestSort(event, property);
  };
  return (
    <TableHead>
      <TableRow>
        <TableCell padding="checkbox">
          <Checkbox
            indeterminate={numSelected > 0 && numSelected < rowCount}
            checked={rowCount > 0 && numSelected === rowCount}
            onChange={onSelectAllClick}
          />
        </TableCell>
        <TableCell padding="checkbox">
        </TableCell>
        {headCells.map((headCell) => (
          <TableCell
            key={headCell.id}
            align={headCell.numeric ? 'right' : 'left'}
            padding={headCell.disablePadding ? 'none' : 'default'}
            sortDirection={orderBy === headCell.id ? order : false}
          >
            <TableSortLabel
              active={orderBy === headCell.id}
              direction={orderBy === headCell.id ? order : 'asc'}
              onClick={createSortHandler(headCell.id)}
            >
              {headCell.label}
              {orderBy === headCell.id ? (
                <span className={classes.visuallyHidden}>
                  {order === 'desc' ? 'sorted descending' : 'sorted ascending'}
                </span>
              ) : null}
            </TableSortLabel>
          </TableCell>
        ))}
      </TableRow>
    </TableHead>
  );
}

const EnhancedTableToolbar = (props) => {
  const classes = useToolbarStyles();
  const { numSelected, onDeleteClick } = props;
  return (
    <Toolbar
      className={clsx(classes.root, {
        [classes.highlight]: numSelected > 0,
      })}>
      {numSelected > 0 ? (
        <Typography className={classes.title} color="inherit" variant="subtitle1" component="div">
          {numSelected} selected
        </Typography>
      ) : (
        <Typography className={classes.title} variant="h6" id="tableTitle" component="div">
          Clusters
        </Typography>
      )}
      <IconButton disabled={numSelected > 0?false:true} onClick={() => onDeleteClick()}>
        <DeleteIcon />
      </IconButton>
    </Toolbar>
  );
};
  
export default function ClusterTable({ rows, onDeleteClusters }) {
    const classes = useStyles();
    
    const [order, setOrder] = React.useState('asc');
    const [orderBy, setOrderBy] = React.useState('calories');
    const [selected, setSelected] = React.useState([]);

    const isSelected = (name) => selected.indexOf(name) !== -1;

    const handleRequestSort = (event, property) => {
      const isAsc = orderBy === property && order === 'asc';
      setOrder(isAsc ? 'desc' : 'asc');
      setOrderBy(property);
    };
  
    const handleSelectAllClick = (event) => {
      if (event.target.checked) {
        const newSelecteds = rows.map((n) => n.ID);
        setSelected(newSelecteds);
        return;
      }
      setSelected([]);
    };

    const handleClick = (event, id) => {
      const selectedIndex = selected.indexOf(id);
      let newSelected = [];
  
      if (selectedIndex === -1) {
        newSelected = newSelected.concat(selected, id);
      } else if (selectedIndex === 0) {
        newSelected = newSelected.concat(selected.slice(1));
      } else if (selectedIndex === selected.length - 1) {
        newSelected = newSelected.concat(selected.slice(0, -1));
      } else if (selectedIndex > 0) {
        newSelected = newSelected.concat(
          selected.slice(0, selectedIndex),
          selected.slice(selectedIndex + 1),
        );
      }
      setSelected(newSelected);
    };

    const descendingComparator = (a, b, orderBy) => {
      if (b[orderBy] < a[orderBy]) {
        return -1;
      }
      if (b[orderBy] > a[orderBy]) {
        return 1;
      }
      return 0;
    }
        
    const getComparator = (order, orderBy) => {
      return order === 'desc'
        ? (a, b) => descendingComparator(a, b, orderBy)
        : (a, b) => -descendingComparator(a, b, orderBy);
    }
        
    const stableSort = (array, comparator) => {
      const stabilizedThis = array.map((el, index) => [el, index]);
      stabilizedThis.sort((a, b) => {
        const order = comparator(a[0], b[0]);
        if (order !== 0) return order;
        return a[1] - b[1];
      });
      return stabilizedThis.map((el) => el[0]);
    }

    return (
      <div className={classes.container}>
        <Paper className={classes.paper}>
          <EnhancedTableToolbar numSelected={selected.length} onDeleteClick={() => onDeleteClusters(selected)}/>
          <TableContainer>
            <Table stickyHeader className={classes.table}>
              <EnhancedTableHead
                classes={classes}
                numSelected={selected.length}
                order={order}
                orderBy={orderBy}
                onSelectAllClick={handleSelectAllClick}
                onRequestSort={handleRequestSort}
                rowCount={rows.length}
              />
              <TableBody>
                {stableSort(rows, getComparator(order, orderBy))
                  .map((row) => {
                    return (
                      <Row key={row.ID} row={row} selected={isSelected(row.ID)} onSelect={(event) => handleClick(event, row.ID)} />
                    );
                  })
                }
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      </div>
    );
}
