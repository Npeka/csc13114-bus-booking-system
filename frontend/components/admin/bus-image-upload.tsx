"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Upload, X, ImageIcon, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import { uploadBusImages, deleteBusImage } from "@/lib/api/trip/bus-service";

interface BusImageUploadProps {
  busId: string;
  imageUrls?: string[];
}

export function BusImageUpload({ busId, imageUrls = [] }: BusImageUploadProps) {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);
  const queryClient = useQueryClient();

  const uploadMutation = useMutation({
    mutationFn: (files: File[]) => uploadBusImages(busId, files),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bus", busId] });
      queryClient.invalidateQueries({ queryKey: ["admin-buses"] });
      setSelectedFiles([]);
      toast.success("Đã tải ảnh lên thành công");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể tải ảnh lên");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (imageUrl: string) => deleteBusImage(busId, imageUrl),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bus", busId] });
      queryClient.invalidateQueries({ queryKey: ["admin-buses"] });
      toast.success("Đã xóa ảnh");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Không thể xóa ảnh");
    },
  });

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);

    // Validate total images
    const totalImages = imageUrls.length + selectedFiles.length + files.length;
    if (totalImages > 10) {
      toast.error("Mỗi xe chỉ được tối đa 10 ảnh");
      return;
    }

    // Validate each file
    const validFiles: File[] = [];
    for (const file of files) {
      if (!file.type.startsWith("image/")) {
        toast.error(`${file.name} không phải là file ảnh`);
        continue;
      }
      if (file.size > 5 * 1024 * 1024) {
        toast.error(`${file.name} vượt quá 5MB`);
        continue;
      }
      validFiles.push(file);
    }

    setSelectedFiles((prev) => [...prev, ...validFiles]);
  };

  const handleUpload = () => {
    if (selectedFiles.length === 0) {
      toast.error("Vui lòng chọn ảnh");
      return;
    }
    uploadMutation.mutate(selectedFiles);
  };

  const handleRemoveSelected = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const handleDeleteImage = (imageUrl: string) => {
    if (confirm("Bạn có chắc chắn muốn xóa ảnh này?")) {
      deleteMutation.mutate(imageUrl);
    }
  };

  const isLoading = uploadMutation.isPending || deleteMutation.isPending;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <ImageIcon className="h-5 w-5" />
          Hình ảnh xe ({imageUrls.length}/10)
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Upload Section */}
        <div className="space-y-3">
          <div className="flex gap-2">
            <label className="flex-1">
              <input
                type="file"
                accept="image/*"
                multiple
                onChange={handleFileSelect}
                className="hidden"
                disabled={isLoading || imageUrls.length >= 10}
              />
              <Button
                type="button"
                variant="outline"
                className="w-full"
                disabled={isLoading || imageUrls.length >= 10}
                asChild
              >
                <span>
                  <Upload className="mr-2 h-4 w-4" />
                  Chọn ảnh
                </span>
              </Button>
            </label>
            {selectedFiles.length > 0 && (
              <Button
                type="button"
                onClick={handleUpload}
                disabled={isLoading}
                className="px-8"
              >
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Đang tải...
                  </>
                ) : (
                  `Tải lên (${selectedFiles.length})`
                )}
              </Button>
            )}
          </div>

          {selectedFiles.length > 0 && (
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4">
              {selectedFiles.map((file, index) => (
                <div key={index} className="group relative">
                  <img
                    src={URL.createObjectURL(file)}
                    alt={`Preview ${index + 1}`}
                    className="h-24 w-full rounded-lg object-cover"
                  />
                  <Button
                    type="button"
                    size="icon"
                    variant="destructive"
                    className="absolute -top-2 -right-2 h-6 w-6 rounded-full opacity-0 group-hover:opacity-100"
                    onClick={() => handleRemoveSelected(index)}
                    disabled={isLoading}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                  <div className="mt-1 truncate text-xs text-muted-foreground">
                    {file.name}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Existing Images */}
        {imageUrls.length > 0 && (
          <div className="space-y-2 border-t pt-4">
            <h4 className="text-sm font-medium">Ảnh hiện có</h4>
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4">
              {imageUrls.map((url, index) => (
                <div key={url} className="group relative">
                  <img
                    src={url}
                    alt={`Bus image ${index + 1}`}
                    className="h-24 w-full rounded-lg object-cover"
                  />
                  <Button
                    type="button"
                    size="icon"
                    variant="destructive"
                    className="absolute -top-2 -right-2 h-6 w-6 rounded-full opacity-0 group-hover:opacity-100"
                    onClick={() => handleDeleteImage(url)}
                    disabled={isLoading}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
              ))}
            </div>
          </div>
        )}

        {imageUrls.length === 0 && selectedFiles.length === 0 && (
          <div className="rounded-lg border border-dashed p-8 text-center">
            <ImageIcon className="mx-auto h-12 w-12 text-muted-foreground/50" />
            <p className="mt-2 text-sm text-muted-foreground">
              Chưa có ảnh nào. Chọn ảnh để tải lên.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
