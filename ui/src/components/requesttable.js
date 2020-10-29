import React from 'react';
import clsx from 'clsx';

import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import IconButton from '@material-ui/core/IconButton';
import CloudDownloadIcon from '@material-ui/icons/CloudDownload';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import Toolbar from '@material-ui/core/Toolbar';
import Checkbox from '@material-ui/core/Checkbox';
import Box from '@material-ui/core/Box';
import Typography from '@material-ui/core/Typography';
import Collapse from '@material-ui/core/Collapse';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
  
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
  { id: 'RequestDate', numeric: false, disablePadding: false, label: 'Created' },
  { id: 'Requested', numeric: true, disablePadding: false, label: '# of Clusters' },
  { id: 'RequestedBy', numeric: false, disablePadding: false, label: 'Requested by' },
  { id: 'DeleteInHours', numeric: true, disablePadding: false, label: 'Delete in hours' },
  { id: 'Status', numeric: false, disablePadding: false, label: 'Status' },
];

function Row(props) {
  const { row, onSelect, onExport, selected } = props;
  const [open, setOpen] = React.useState(false);
  const classes = useRowStyles();
  let rowDate = new Date(0);
  rowDate.setUTCSeconds(row.Created);
  let rowDateStr = (rowDate.getMonth()+1) + "-" + rowDate.getDate() + "-" + rowDate.getFullYear();
  return (
    <React.Fragment>
      <TableRow className={classes.root} key={row.ID} hover onClick={() => onSelect(row)} aria-checked={selected} selected={selected}>
        <TableCell padding="checkbox">
          <Checkbox
            checked={selected}
            onClick={() => onSelect(row)}
          />
        </TableCell>
        <TableCell>
          <IconButton aria-label="expand row" size="small" onClick={(event) => {event.stopPropagation(); setOpen(!open);}}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row" className={classes.oneLine}>{row.ID.substring(row.ID.length - 5, row.ID.length)}</TableCell>
        <TableCell align="left">{rowDateStr}</TableCell>
        <TableCell align="right">{row.Requested}</TableCell>
        <TableCell align="left">{row.RequestedBy}</TableCell>
        <TableCell align="right">{row.DeleteInHours}</TableCell>
        <TableCell align="left">{row.Status}</TableCell>
        <TableCell>
          <IconButton aria-label='export' color='primary' onClick={() => onExport(row)}>
            <CloudDownloadIcon/>
          </IconButton>
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={9}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box margin={1}>
              <Typography variant="h6" gutterBottom component="div">Request Details</Typography>
              <Table>
                  <tbody>
                      <tr><td><Typography>Request Id:</Typography></td><td>{row.ID}</td></tr>
                      <tr><td><Typography>Number of Clusters:</Typography></td><td>{row.Requested}</td></tr>
                      <tr><td><Typography>No Subnet:</Typography></td><td>{row.NoSubnet?'true':'false'}</td></tr>
                      <tr><td><Typography>Zone:</Typography></td><td>{row.Zone}</td></tr>
                      <tr><td><Typography>Requested at:</Typography></td><td>{rowDate.toString()}</td></tr>
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
        <TableCell padding="checkbox">
        </TableCell>
      </TableRow>
    </TableHead>
  );
}

const EnhancedTableToolbar = (props) => {
  const classes = useToolbarStyles();
  const { numSelected } = props;
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
          Cluster Requests
        </Typography>
      )}
    </Toolbar>
  );
};

export default function RequestTable({ rows, onSelect, onExport }) {
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
    if (selected[0] && selected[0]===id) {
      setSelected([]);
    } else {
      setSelected([id]);
      rows.forEach(row => {
        if (row.ID === id)
          onSelect(row);
      });
    }
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
        <EnhancedTableToolbar numSelected={selected.length}/>
        <TableContainer component={Paper} style={{ overflow: "auto" }}>
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
                    <Row key={row.ID} row={row} selected={isSelected(row.ID)} onSelect={(event) => handleClick(event, row.ID)} onExport={onExport}/>
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
