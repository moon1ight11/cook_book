
export type Role = 'Basic' | 'Admin';
export type Difficulty = 'easy' | 'medium' | 'hard';

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  created_at: string;
  updated_at: string;
}

export interface Ingredient {
  name: string;
  quantity: string;
  unit: string;
}

export interface Recipe {
  id: string;
  owner_id: string;
  title: string;
  description: string;
  prep_time: number;
  cook_time: number;
  servings: number;
  difficulty: Difficulty | '';
  tags: string[];
  created_at: string;
}

export interface DetailedRecipe extends Recipe {
  ingredients: Ingredient[];
  steps: string[];
  updated_at: string;
}

export interface Paginated<T> {
  data: T[] | null;
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ErrorResponse {
  error: string;
}

export interface SuccessResponse {
  message: string;
}

export interface RefreshTokenResponse {
  refresh_token: string;
}

// --- запросы ---

export interface SignUpRequest {
  name: string;
  email: string;
  password: string;
}

export interface SignInRequest {
  email: string;
  password: string;
}

export interface CreateRecipeRequest {
  title: string;
  description?: string;
  ingredients: Ingredient[];
  steps: string[];
  prep_time?: number;
  cook_time?: number;
  servings: number;
  difficulty?: Difficulty;
  tags?: string[];
}

export type UpdateRecipeRequest = Partial<CreateRecipeRequest>;

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  password?: string;
}

export interface ListRecipesParams {
  page?: number;
  page_size?: number;
  title?: string;
  difficulty?: Difficulty;
}

export interface CreateAdminRequest {
  name: string;
  email: string;
  password: string;
}

export type RecipeList = Recipe[] | null;