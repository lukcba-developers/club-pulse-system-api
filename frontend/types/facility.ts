export enum FacilityType {
    Court = "court",
    Pool = "pool",
    Gym = "gym",
    Field = "field",
    // Legacy values that may exist in DB
    TennisCourt = "Tennis Court",
    PadelCourt = "Padel Court",
    SwimmingPool = "Swimming Pool",
    FootballField = "Football Field",
    GolfSimulator = "Golf Simulator"
}

export enum FacilityStatus {
    Active = "active",
    Maintenance = "maintenance",
    Closed = "closed"
}

export interface FacilitySpecifications {
    surface_type?: string;
    lighting: boolean;
    covered: boolean;
    equipment?: string[];
}

export interface FacilityLocation {
    name: string;
    description?: string;
}

export interface Facility {
    id: string;
    club_id: string;
    name: string;
    description?: string;
    type: FacilityType;
    status: FacilityStatus;
    capacity: number;
    hourly_rate: number;
    opening_time: string; // HH:MM
    closing_time: string; // HH:MM
    guest_fee: number;
    specifications: FacilitySpecifications;
    location: FacilityLocation;
    created_at: string;
    updated_at: string;
}
