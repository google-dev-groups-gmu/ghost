export interface ScheduleLabel {
    id: string;
    course_id: string;
    section: string;
    professor: string;
}

export interface ScheduleItem {
    day: number;
    start_time: number;
    end_time: number;
    location: string;
    label: ScheduleLabel[];
}

export interface Room {
    ID: string;
    Building: string;
    Number: string;
    Schedule: ScheduleItem[];
}

export interface BuildingDrawerProps {
    buildingName: string;
    buildingId: string | null;
    onClose: () => void;
}

export const DAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
