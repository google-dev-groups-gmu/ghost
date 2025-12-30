"use client";

import { useState } from "react";
import Map, { Layer, Marker, Source } from "react-map-gl/mapbox";
import "mapbox-gl/dist/mapbox-gl.css";
import { buildings, type buildingInfo } from "@/types/buildings";
import { campusMask, MAPBOX_TOKEN, mapStyle, MAX_BOUNDS } from "@/types/map";
import { Dot } from "lucide-react";
import { BuildingDrawer } from "@/components/drawer/drawer";

export default function CampusMap() {
    const [hoverInfo, setHoverInfo] = useState<buildingInfo | null>(null);
    const [selectedBuildingId, setSelectedBuildingId] = useState<string | null>(
        null
    );
    const [selectedBuildingName, setSelectedBuildingName] =
        useState<string>("");

    return (
        <div className="relative w-screen h-screen">
            <Map
                initialViewState={{
                    zoom: 16,
                    pitch: 30,
                    bearing: -20,
                    longitude: -77.30761744755588,
                    latitude: 38.83006053754113,
                }}
                mapStyle={mapStyle}
                mapboxAccessToken={MAPBOX_TOKEN}
                maxBounds={MAX_BOUNDS}
                dragRotate={false}
                onLoad={(e) => {
                    const map = e.target;

                    try {
                        map.setConfigProperty(
                            "basemap",
                            "showPointOfInterestLabels",
                            false
                        );
                        map.setConfigProperty(
                            "basemap",
                            "showPlaceLabels",
                            false
                        );
                        map.setConfigProperty(
                            "basemap",
                            "showRoadLabels",
                            false
                        );
                        map.setConfigProperty(
                            "basemap",
                            "showTransitLabels",
                            false
                        );
                    } catch (error) {}

                    const layersToHide = [
                        "poi-label",
                        "road-label",
                        "transit-label",
                    ];
                    layersToHide.forEach((layer) => {
                        if (map.getLayer(layer)) {
                            map.setLayoutProperty(layer, "visibility", "none");
                        }
                    });
                }}
            >
                <Source id="mask-source" type="geojson" data={campusMask}>
                    <Layer
                        id="world-mask"
                        type="fill"
                        paint={{
                            "fill-color": "rgba(0, 0, 0, 0.5)",
                        }}
                    />
                </Source>
                {buildings.map((b) => (
                    <div key={b.id}>
                        <Marker
                            key={b.id}
                            longitude={b.lng}
                            latitude={b.lat}
                            anchor="bottom"
                            onClick={(e) => {
                                e.originalEvent.stopPropagation();
                                setSelectedBuildingId(b.id);
                                setSelectedBuildingName(b.name);
                            }}
                        >
                            <div className="flex flex-col items-center">
                                <div
                                    className="w-fit px-1.5 py-0.5 rounded cursor-pointer flex items-center justify-center shadow-lg bg-black/25 text-white text-xs"
                                    onMouseEnter={() => setHoverInfo(b)}
                                    onMouseLeave={() => setHoverInfo(null)}
                                >
                                    {hoverInfo?.id === b.id ? (
                                        <span className="text-sm">
                                            {b.name}
                                        </span>
                                    ) : (
                                        <span>{b.id}</span>
                                    )}
                                </div>
                                <Dot className="text-white" />
                            </div>
                        </Marker>
                    </div>
                ))}
            </Map>

            <BuildingDrawer
                buildingName={selectedBuildingName}
                buildingId={selectedBuildingId}
                onClose={() => setSelectedBuildingId(null)}
            />
        </div>
    );
}
