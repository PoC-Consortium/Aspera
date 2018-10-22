import { AuthApiActions } from '../actions';
import { AccountCreateActionsTypes, AccountCreateActionsUnion } from '../../pages/dashboard/setup/account/create.actions';

export interface State {
  error: string | null;
  pending: boolean;
}

export const initialState: State = {
  error: null,
  pending: false,
};

export function reducer(
  state = initialState,
  action:
    | AuthApiActions.AuthApiActionsUnion
    | AccountCreateActionsUnion
): State {
  switch (action.type) {
    case AccountCreateActionsTypes.CreateActiveAccount: {
      return {
        ...state,
        error: null,
        pending: true,
      };
    }

    case AuthApiActions.AuthApiActionTypes.LoginSuccess: {
      return {
        ...state,
        error: null,
        pending: false,
      };
    }

    case AuthApiActions.AuthApiActionTypes.LoginFailure: {
      return {
        ...state,
        error: action.payload.error,
        pending: false,
      };
    }

    default: {
      return state;
    }
  }
}

export const getError = (state: State) => state.error;
export const getPending = (state: State) => state.pending;
