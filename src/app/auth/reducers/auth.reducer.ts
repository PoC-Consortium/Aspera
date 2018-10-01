import { AuthApiActions, AuthActions } from '../actions';
import { Account } from '../../lib/model/account';

export interface State {
  account: Account | null;
}

export const initialState: State = {
  account: null,
};

export function reducer(
  state = initialState,
  action: AuthApiActions.AuthApiActionsUnion | AuthActions.AuthActionsUnion
): State {
  switch (action.type) {
    case AuthApiActions.AuthApiActionTypes.LoginSuccess: {
      return {
        ...state,
        account: action.payload.account,
      };
    }

    case AuthActions.AuthActionTypes.Logout: {
      return initialState;
    }

    default: {
      return state;
    }
  }
}

export const getAccount = (state: State) => state.account;
