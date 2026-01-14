export enum UserRole {
    SUPER_ADMIN = 'SUPER_ADMIN',
    ADMIN = 'ADMIN',
    MEMBER = 'MEMBER',
    COACH = 'COACH',
    MEDICAL_STAFF = 'MEDICAL_STAFF'
}

export interface UserProfile {
    id: string;
    email: string;
    name: string;
    role: UserRole;
    club_id: string;
    avatar_url?: string;
}
