import React, {Component} from 'react';
import AlbumItem from "./AlbumItem";
import PropTypes from 'prop-types';

class Albums extends Component {
    render() {
        return this.props.albums.map(album => (
            <AlbumItem key={album.id} album={album}/>
        ))
    }
}

Albums.propTypes = {
    albums: PropTypes.array.isRequired
};

export default Albums;
