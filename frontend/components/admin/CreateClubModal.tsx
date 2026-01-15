'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import Image from 'next/image';
import { clubService } from '@/services/club-service';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

const clubSchema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    slug: z.string().optional(),
    domain: z.string().optional(),
    logo_url: z.string().url('Invalid URL').optional().or(z.literal('')),
    primary_color: z.string().regex(/^#[0-9A-F]{6}$/i, 'Invalid Color').optional().or(z.literal('')),
    secondary_color: z.string().regex(/^#[0-9A-F]{6}$/i, 'Invalid Color').optional().or(z.literal('')),
    contact_email: z.string().email('Invalid Email').optional().or(z.literal('')),
    contact_phone: z.string().optional(),
    theme_config: z.string().optional(),
    settings: z.string().optional(),
});

type ClubFormData = z.infer<typeof clubSchema>;

interface CreateClubModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onSuccess: () => void;
}

export function CreateClubModal({ open, onOpenChange, onSuccess }: CreateClubModalProps) {
    const [loading, setLoading] = useState(false);
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const { register, handleSubmit, watch, formState: { errors }, reset } = useForm<ClubFormData>({
        resolver: zodResolver(clubSchema),
        defaultValues: {
            name: '',
            slug: '',
            domain: '',
            logo_url: '',
            primary_color: '#000000',
            secondary_color: '#ffffff',
            contact_email: '',
            contact_phone: '',
            theme_config: '',
            settings: ''
        }
    });

    const watchedName = watch('name');
    const watchedLogo = watch('logo_url');
    const watchedPrimary = watch('primary_color');
    const watchedSecondary = watch('secondary_color');

    // Preview logo: either selected file blob or entered URL
    const logoPreview = selectedFile ? URL.createObjectURL(selectedFile) : watchedLogo;

    const onSubmit = async (data: ClubFormData) => {
        setLoading(true);
        let createdClubId: string | null = null;
        try {
            // 1. Create Club
            const newClub = await clubService.createClub({
                ...data,
                slug: data.slug || undefined,
                domain: data.domain || undefined,
                logo_url: data.logo_url || undefined,
                primary_color: data.primary_color || undefined,
                secondary_color: data.secondary_color || undefined,
                contact_email: data.contact_email || undefined,
                contact_phone: data.contact_phone || undefined,
                theme_config: data.theme_config || undefined,
                settings: data.settings || undefined,
            });
            createdClubId = newClub.id;

            // 2. Upload Logo if selected
            if (selectedFile && newClub.id) {
                try {
                    await clubService.uploadLogo(newClub.id, selectedFile);
                } catch (uploadErr) {
                    console.error("Logo upload failed:", uploadErr);
                    alert("Club created successfully, but logo upload failed. You can try uploading it again from the club settings.");
                    // We still proceed to success because the club exists
                }
            }

            onSuccess();
            onOpenChange(false);
            reset();
            setSelectedFile(null);
        } catch (err: unknown) {
            console.error(err);
            const errorMessage = err instanceof Error ? err.message : 'Failed to create club';
            // Only show main failure if club wasn't created
            if (!createdClubId) {
                alert(errorMessage);
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[700px]">
                <DialogHeader>
                    <DialogTitle>Create New Club</DialogTitle>
                    <DialogDescription>
                        Define the branding and settings for the new club.
                    </DialogDescription>
                </DialogHeader>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <form id="create-club-form" onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                        <div className="grid gap-2">
                            <Label htmlFor="name">Name</Label>
                            <Input id="name" {...register('name')} />
                            {errors.name && <p className="text-red-500 text-xs">{errors.name.message}</p>}
                        </div>
                        <div className="grid gap-2">
                            <Label htmlFor="slug">Slug (Auto-generated if empty)</Label>
                            <Input id="slug" {...register('slug')} placeholder="club-name" />
                        </div>
                        <div className="grid gap-2">
                            <Label htmlFor="domain">Domain (Optional)</Label>
                            <Input id="domain" {...register('domain')} placeholder="club.com" />
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="primary_color">Primary Color</Label>
                                <div className="flex gap-2">
                                    <Input id="primary_color" type="color" className="w-12 p-1" {...register('primary_color')} />
                                    <Input {...register('primary_color')} placeholder="#000000" />
                                </div>
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="secondary_color">Secondary Color</Label>
                                <div className="flex gap-2">
                                    <Input id="secondary_color" type="color" className="w-12 p-1" {...register('secondary_color')} />
                                    <Input {...register('secondary_color')} placeholder="#ffffff" />
                                </div>
                            </div>
                        </div>
                        <div className="grid gap-2">
                            <Label htmlFor="logo_file">Logo</Label>
                            <Input
                                id="logo_file"
                                type="file"
                                accept="image/*"
                                onChange={(e) => {
                                    if (e.target.files && e.target.files[0]) {
                                        setSelectedFile(e.target.files[0]);
                                    }
                                }}
                            />
                            <div className="text-xs text-center text-gray-500">- OR -</div>
                            <Input id="logo_url" {...register('logo_url')} placeholder="https://... (External URL)" />
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="contact_email">Contact Email</Label>
                                <Input id="contact_email" type="email" {...register('contact_email')} />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="contact_phone">Contact Phone</Label>
                                <Input id="contact_phone" {...register('contact_phone')} />
                            </div>
                        </div>
                    </form>

                    {/* Preview Section */}
                    <div className="hidden md:block border rounded-lg p-4 bg-gray-50">
                        <h3 className="text-sm font-semibold mb-3 text-gray-500 uppercase tracking-wider">Preview</h3>
                        <div className="border rounded-lg bg-white shadow-sm overflow-hidden min-h-[300px] flex flex-col font-sans">
                            {/* Header Preview */}
                            <div className="p-4 border-b flex justify-between items-center" style={{ backgroundColor: '#ffffff' }}>
                                <div className="flex items-center gap-2">
                                    {logoPreview ? (
                                        <div className="relative w-8 h-8">
                                            <Image
                                                src={logoPreview}
                                                alt="Logo"
                                                fill
                                                className="object-contain"
                                                unoptimized
                                                onError={(e) => {
                                                    // Hide parent or handle error? Generic fallback handled by parent logic usually.
                                                    // next/image doesn't have onError in same way on div.
                                                    // We can use style display none on the wrapper if image fails?
                                                    // For preview, simple img might be better if Image proves tricky with blob: URLs.
                                                    // BUT user asked to fix warning. "unoptimized" works with blobs usually.
                                                    const target = e.target as HTMLImageElement;
                                                    target.style.display = 'none';
                                                }}
                                            />
                                        </div>
                                    ) : (
                                        <div className="h-8 w-8 bg-gray-200 rounded-full"></div>
                                    )}
                                    <span className="font-bold text-lg text-gray-800">{watchedName || 'Club Name'}</span>
                                </div>
                                <div className="flex gap-2 text-sm">
                                    <span className="text-gray-600">Home</span>
                                    <span className="text-gray-600">Bookings</span>
                                </div>
                            </div>
                            {/* Body Preview */}
                            <div className="p-6 flex-1 bg-white">
                                <h1 className="text-2xl font-bold mb-4" style={{ color: watchedPrimary }}>Welcome to {watchedName || 'Club Name'}</h1>
                                <button className="px-4 py-2 rounded text-white font-medium" style={{ backgroundColor: watchedPrimary }}>
                                    Book Now
                                </button>
                                <div className="mt-4 p-4 border rounded" style={{ borderColor: watchedSecondary, borderLeftWidth: '4px' }}>
                                    <h4 className="font-semibold text-gray-700">Next Match</h4>
                                    <p className="text-sm text-gray-500">Tomorrow at 19:00</p>
                                </div>
                            </div>
                        </div>
                        <p className="text-xs text-center text-gray-400 mt-2">This is a preview of how the club branding will look.</p>
                    </div>
                </div>
                <DialogFooter className="mt-4">
                    <Button type="submit" form="create-club-form" disabled={loading}>
                        {loading ? 'Creating...' : 'Create Club'}
                    </Button>
                </DialogFooter>
            </DialogContent >
        </Dialog >
    );
}
