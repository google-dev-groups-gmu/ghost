"use client";

import { useState, useEffect } from "react";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Room } from "@/types/drawer";
import { Loader2, X } from "lucide-react";
import { BACKEND_URL } from "@/lib/utils";
import { BuildingDrawerProps } from "@/types/drawer";

const formatTime = (minutes: number) => {
    const h = Math.floor(minutes / 60);
    const m = minutes % 60;
    const ampm = h >= 12 ? "PM" : "AM";
    const hour = h % 12 || 12;
    return `${hour}:${m.toString().padStart(2, "0")} ${ampm}`;
};

export function BuildingDrawer({
    buildingName,
    buildingId,
    onClose,
}: BuildingDrawerProps) {
    const [rooms, setRooms] = useState<Room[]>([]);
    const [loading, setLoading] = useState(false);

    const now = new Date();
    const currentMinutes = now.getHours() * 60 + now.getMinutes();

    useEffect(() => {
        if (!buildingId) {
            setRooms([]);
            return;
        }

        const fetchData = async () => {
            setLoading(true);
            try {
                const now = new Date();
                const currentMinutes = now.getHours() * 60 + now.getMinutes();
                const today = now.getDay();

                const res = await fetch(
                    `${BACKEND_URL}/api/rooms?building=${buildingId}&day=${today}&time=${currentMinutes}`
                );

                if (res.ok) {
                    const data = await res.json();
                    setRooms(data || []);
                } else {
                    console.error("server error:", res.statusText);
                }
            } catch (error) {
                console.error("failed to fetch rooms", error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [buildingId]);

    const roomStatuses = (rooms || [])
        .flatMap(
            (
                room
            ): {
                roomNumber: string;
                isOccupied: boolean;
                startTime: number | null;
                endTime: number | null;
                course_id: string | null;
                section: string | null;
                professor: string | null;
            }[] => {
                const activeSchedules = room.Schedule || [];

                if (activeSchedules.length > 0) {
                    return activeSchedules.flatMap((schedule) =>
                        schedule.label.map((cls) => ({
                            roomNumber: room.Number,
                            isOccupied: true,
                            startTime: schedule.start_time,
                            endTime: schedule.end_time,
                            course_id: cls.course_id,
                            section: cls.section,
                            professor: cls.professor,
                        }))
                    );
                }

                return [
                    {
                        roomNumber: room.Number,
                        isOccupied: false,
                        startTime: null,
                        endTime: null,
                        course_id: null,
                        section: null,
                        professor: null,
                    },
                ];
            }
        )
        .sort((a, b) =>
            a.roomNumber.localeCompare(b.roomNumber, undefined, {
                numeric: true,
            })
        );

    if (!buildingId) return null;

    return (
        <div className="fixed inset-0 z-50 flex justify-end">
            <div
                className="fixed inset-0 bg-black/50 transition-opacity"
                onClick={onClose}
            />

            <div className="relative z-50 w-full lg:max-w-xl h-full bg-background shadow-xl border-l flex flex-col animate-in slide-in-from-right duration-300">
                <div className="flex items-center justify-between p-4 border-b">
                    <div className="space-y-1">
                        <span className="text-lg font-semibold">
                            {buildingName}
                        </span>
                        <p className="text-sm text-muted-foreground">
                            Current room availability (
                            {formatTime(currentMinutes)})
                        </p>
                    </div>
                    <button
                        onClick={onClose}
                        className="rounded-sm opacity-70 hover:opacity-100"
                    >
                        <X className="h-4 w-4" />
                        <span className="sr-only">Close</span>
                    </button>
                </div>

                <div className="flex-1 overflow-y-auto px-4 pb-4">
                    {loading ? (
                        <div className="flex h-40 items-center justify-center">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>Room</TableHead>
                                    <TableHead>Time</TableHead>
                                    <TableHead>Status / Course</TableHead>
                                    <TableHead>Professor</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {roomStatuses.length > 0 ? (
                                    roomStatuses.map((status, idx) => (
                                        <TableRow
                                            key={`${status.roomNumber}-${idx}`}
                                        >
                                            <TableCell>
                                                {status.roomNumber}
                                            </TableCell>
                                            <TableCell className="text-xs text-muted-foreground whitespace-nowrap">
                                                {status.isOccupied
                                                    ? `${formatTime(
                                                          status.startTime!
                                                      )} - ${formatTime(
                                                          status.endTime!
                                                      )}`
                                                    : "-"}
                                            </TableCell>
                                            <TableCell>
                                                {status.isOccupied ? (
                                                    <>
                                                        <div>
                                                            {status.course_id}
                                                        </div>
                                                        <div className="text-xs text-muted-foreground">
                                                            Sec {status.section}
                                                        </div>
                                                    </>
                                                ) : (
                                                    <span className="inline-flex items-center rounded-md px-1.5 py-0.5 text-xs ring-1 ring-inset">
                                                        Empty
                                                    </span>
                                                )}
                                            </TableCell>
                                            <TableCell className="text-sm">
                                                {status.isOccupied
                                                    ? status.professor
                                                    : "-"}
                                            </TableCell>
                                        </TableRow>
                                    ))
                                ) : (
                                    <TableRow>
                                        <TableCell
                                            colSpan={4}
                                            className="text-center h-24 text-muted-foreground"
                                        >
                                            No rooms found.
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    )}
                </div>
            </div>
        </div>
    );
}
