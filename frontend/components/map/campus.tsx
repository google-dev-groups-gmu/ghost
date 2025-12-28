"use client";

import { useState } from "react";
import Map, { Marker, Popup } from "react-map-gl/mapbox";
import "mapbox-gl/dist/mapbox-gl.css";
import { useRouter } from "next/navigation";
import { buildings, type buildingInfo } from "@/types/buildings";
import { MAPBOX_TOKEN, mapStyle, MAX_BOUNDS } from "@/types/map";

export default function CampusMap() {
    const router = useRouter();
    const [hoverInfo, setHoverInfo] = useState<buildingInfo | null>(null);

    return (
        <div className="relative w-screen h-screen bg-[#1a1a1a]">
            <Map
                initialViewState={{
                    zoom: 15,
                    pitch: 30,
                    bearing: -20,
                }}
                style={{ width: "100vw", height: "100vh" }}
                mapStyle={mapStyle}
                mapboxAccessToken={MAPBOX_TOKEN}
                maxBounds={MAX_BOUNDS}
                dragRotate={false}
            >
                {buildings.map((b) => (
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
                        <div
                            className="w-3 h-3 rounded-full cursor-pointer hover:scale-150 transition-transform"
                            onMouseEnter={() => setHoverInfo(b)}
                            onMouseLeave={() => setHoverInfo(null)}
                            style={{
                                backgroundColor:
                                    b.status === "ghost"
                                        ? "#ffffffff"
                                        : "#000000ff",
                            }}
                        />
                    </Marker>
                ))}

                {hoverInfo && (
                    <Popup
                        longitude={hoverInfo.lng}
                        latitude={hoverInfo.lat}
                        offset={20}
                        closeButton={false}
                        closeOnClick={false}
                        anchor="right"
                    >
                        <div className="text-black font-sans text-center mt-1 mx-1">
                            {hoverInfo.name}
                        </div>
                    </Popup>
                )}
            </Map>
        </div>
    );
}
