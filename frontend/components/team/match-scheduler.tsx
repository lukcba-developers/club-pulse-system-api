"use client";

import { useState } from "react";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { teamService } from "@/services/team-service";
import { useToast } from "@/components/ui/use-toast";
import { Users } from 'lucide-react';

// Mock training groups for now, or fetch if available
const MOCK_GROUPS = [
    { id: "1", name: "Equipo Titular" },
    { id: "2", name: "Reserva Sub-20" },
    { id: "3", name: "Veteranos +35" }
];

export function MatchScheduler() {
    const [open, setOpen] = useState(false);
    const [loading, setLoading] = useState(false);
    const { toast } = useToast();

    const [form, setForm] = useState({
        training_group_id: "",
        opponent_name: "",
        is_home_game: true,
        meetup_time: "",
        location: ""
    });

    const handleSubmit = async () => {
        try {
            setLoading(true);
            await teamService.scheduleMatch({
                training_group_id: form.training_group_id || "mock-group-id", // Use mock if not selected
                opponent_name: form.opponent_name,
                is_home_game: form.is_home_game,
                meetup_time: new Date(form.meetup_time).toISOString(),
                location: form.location
            });
            toast({
                title: "Partido Programado",
                description: `Partido contra ${form.opponent_name} creado exitosamente.`,
            });
            setOpen(false);
            setForm({ training_group_id: "", opponent_name: "", is_home_game: true, meetup_time: "", location: "" });
        } catch (error) {
            console.error(error);
            toast({
                title: "Error",
                description: "No se pudo programar el partido.",
                variant: "destructive"
            });
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button variant="outline" className="w-full justify-start">
                    <Users className="mr-2 h-4 w-4" />
                    Organizar Partido
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Programar Nuevo Partido</DialogTitle>
                    <DialogDescription>
                        Crea un evento de partido para tu equipo. Los jugadores recibirán una notificación.
                    </DialogDescription>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="group" className="text-right">Equipo</Label>
                        <Select onValueChange={(v) => setForm({ ...form, training_group_id: v })}>
                            <SelectTrigger className="col-span-3">
                                <SelectValue placeholder="Seleccionar Equipo" />
                            </SelectTrigger>
                            <SelectContent>
                                {MOCK_GROUPS.map(g => (
                                    <SelectItem key={g.id} value={g.id}>{g.name}</SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="opponent" className="text-right">Rival</Label>
                        <Input
                            id="opponent"
                            className="col-span-3"
                            value={form.opponent_name}
                            onChange={(e) => setForm({ ...form, opponent_name: e.target.value })}
                        />
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="datetime" className="text-right">Fecha/Hora</Label>
                        <Input
                            id="datetime"
                            type="datetime-local"
                            className="col-span-3"
                            value={form.meetup_time}
                            onChange={(e) => setForm({ ...form, meetup_time: e.target.value })}
                        />
                    </div>
                    <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="location" className="text-right">Ubicación</Label>
                        <Input
                            id="location"
                            className="col-span-3"
                            placeholder="Dirección o Cancha"
                            value={form.location}
                            onChange={(e) => setForm({ ...form, location: e.target.value })}
                        />
                    </div>
                </div>
                <DialogFooter>
                    <Button type="submit" onClick={handleSubmit} disabled={loading}>
                        {loading ? "Creando..." : "Programar Partido"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
