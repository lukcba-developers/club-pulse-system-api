import api from "@/lib/axios";

export interface ScheduleMatchRequest {
    training_group_id: string;
    opponent_name: string;
    is_home_game: boolean;
    meetup_time: string; // ISO string
    location: string;
}

export const teamService = {
    scheduleMatch: async (data: ScheduleMatchRequest) => {
        const response = await api.post("/team/events", data);
        return response.data;
    },

    respondAvailability: async (eventID: string, status: 'CONFIRMED' | 'DECLINED' | 'MAYBE', reason?: string) => {
        const response = await api.post("/team/availability", {
            event_id: eventID,
            status,
            reason
        });
        return response.data;
    }
};
