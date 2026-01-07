import api from '../lib/axios';

export interface AttendanceRecord {
    user_id: string;
    status: 'PRESENT' | 'ABSENT' | 'LATE' | 'EXCUSED';
    has_debt?: boolean;
    user?: {
        name: string;
        email?: string;
    }
}

export interface AttendanceList {
    id: string;
    group: string;
    training_group_id?: string;
    date: string;
    records: AttendanceRecord[];
}

export const attendanceService = {
    getGroupList: async (group: string, date?: string) => {
        const query = date ? `?date=${date}` : '';
        const response = await api.get<AttendanceList>(`/attendance/groups/${group}${query}`);
        return response.data;
    },

    markAttendance: async (listId: string, userId: string, status: string, notes: string = '') => {
        await api.post(`/attendance/${listId}/records`, {
            user_id: userId,
            status,
            notes
        });
    },

    getTrainingGroupList: async (groupId: string, groupName: string, category: string, date?: string) => {
        const query = date ? `&date=${date}` : '';
        const groupNameEncoded = encodeURIComponent(groupName);
        const response = await api.get<AttendanceList>(`/attendance/training-groups/${groupId}?group_name=${groupNameEncoded}&category=${category}${query}`);
        return response.data;
    }
}
