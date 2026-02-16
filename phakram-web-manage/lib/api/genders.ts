import { apiClient } from './client';

export interface Gender {
  id: string;
  name_th: string;
  name_en: string;
  is_active: boolean;
  created_at: string;
}

export interface CreateGenderDto {
  name_th: string;
  name_en: string;
  is_active: boolean;
}

export interface UpdateGenderDto {
  name_th: string;
  name_en: string;
  is_active: boolean;
}

export interface ListGendersParams {
  page?: number;
  limit?: number;
}

export interface ListGendersResponse {
  data: Gender[];
  paginate?: {
    page: number;
    size: number;
    total: number;
  };
}

export const gendersApi = {
  /**
   * Get list of genders with pagination
   */
  async list(params?: ListGendersParams): Promise<ListGendersResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const query = queryParams.toString();
    const endpoint = `/system/genders/${query ? `?${query}` : ''}`;

    const response = await apiClient.get<Gender[]>(endpoint);
    return {
      data: response.data,
      paginate: response.paginate,
    };
  },

  /**
   * Get single gender by ID
   */
  async getById(id: string): Promise<Gender> {
    const response = await apiClient.get<Gender>(`/system/genders/${id}`);
    return response.data;
  },

  /**
   * Create new gender
   */
  async create(data: CreateGenderDto): Promise<void> {
    await apiClient.post('/system/genders/', data);
  },

  /**
   * Update existing gender
   */
  async update(id: string, data: UpdateGenderDto): Promise<void> {
    await apiClient.patch(`/system/genders/${id}`, data);
  },

  /**
   * Delete gender
   */
  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/genders/${id}`);
  },
};
