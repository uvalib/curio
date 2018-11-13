define(function () {
    return function (formats) {
        return {
            async: ['TreeComponent.js', 'GalleryComponent.js', 'MetadataComponent.js', 'openseadragon.min.js']
            //async: ['TreeComponent', 'iiifgallery.proxy', 'GalleryComponent', 'MetadataComponent', 'openseadragon.min']
        };
    };
});
