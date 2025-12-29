"use client";

import { useState } from "react";
import Map, { Layer, Marker, Popup, Source } from "react-map-gl/mapbox";
import "mapbox-gl/dist/mapbox-gl.css";
import { useRouter } from "next/navigation";
import { buildings, type buildingInfo } from "@/types/buildings";
import { campusMask, MAPBOX_TOKEN, mapStyle, MAX_BOUNDS } from "@/types/map";
import { Dot } from "lucide-react";

export default function CampusMap() {
    const router = useRouter();
    const [hoverInfo, setHoverInfo] = useState<buildingInfo | null>(null);

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
                                router.push(`/building/${b.id}`);
                            }}
                        >
                            <div className="flex flex-col items-center">
                                <div
                                    className="w-fit px-2 rounded-md cursor-pointer flex items-center justify-center shadow-lg bg-black/25"
                                    onMouseEnter={() => setHoverInfo(b)}
                                    onMouseLeave={() => setHoverInfo(null)}
                                >
                                    {hoverInfo?.id === b.id ? b.name : b.id}
                                </div>
                                <Dot />
                            </div>
                        </Marker>
                    </div>
                ))}
            </Map>
        </div>
    );
}
