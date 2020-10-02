import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import FormGroup from '@material-ui/core/FormGroup';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Slider from '@material-ui/core/Slider';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Checkbox from '@material-ui/core/Checkbox';
import Button from '@material-ui/core/Button';

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 220,
  },
}));

export default function RequestForm({ zones, onSubmit }) {
  const classes = useStyles();
  
  const [numberOfClusters, setNumberOfClusters] = React.useState(0);
  const [deleteInHours, setDeleteInHours] = React.useState(0);
  const [zone, setZone] = React.useState('');
  const [subnet, setSubnet] = React.useState(false);

  const onClickRequest = () => {
    onSubmit({
        numberOfClusters: numberOfClusters,
        zone: zone,
        deleteInHours: deleteInHours,
        subnet: subnet,
    })
  }

  return (
    <FormGroup row>
        <FormControl className={classes.formControl}>
            <InputLabel id="clustern-label">Number of Clusters</InputLabel>
            <Slider
                defaultValue={10}
                aria-labelledby="clustern-label"
                valueLabelDisplay="auto"
                step={10}
                marks
                min={10}
                max={500}
                onChange={(event, newValue) => setNumberOfClusters(newValue)}
            />
        </FormControl>
        <FormControl className={classes.formControl}>
            <InputLabel id="zone-label">Zone</InputLabel>
            <Select labelId="zone-label" id="zone-select" value={zone} onChange={(event) => setZone(event.target.value)}>
                {zones.map((zone, index) =>
                <MenuItem key={index} value={index}>{zone}</MenuItem>
                )}
            </Select>
        </FormControl>
        <FormControl className={classes.formControl}>
            <InputLabel id="ttl-label">Delete in Hours</InputLabel>
            <Slider
                defaultValue={10}
                aria-labelledby="ttl-label"
                valueLabelDisplay="auto"
                step={24}
                marks
                min={24}
                max={720}
                onChange={(event, newValue) => setDeleteInHours(newValue)}
            />
        </FormControl>
        <FormControl className={classes.formControl}>
            <FormControlLabel
                control={<Checkbox checked={subnet} onChange={(event) => setSubnet(event.target.checked)} name="subnet" />}
                label="Subnet"
            />
        </FormControl>
        <FormControl className={classes.formControl}>
            <Button variant="contained" onClick={() => onClickRequest()}>Request Clusters</Button>
        </FormControl>
    </FormGroup>
  )
};