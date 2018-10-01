import { reducer } from '../reducers/auth.reducer';
import * as fromAuth from '../reducers/auth.reducer';
import { AuthApiActions, AuthActions } from '../actions/';

import { Account } from '../../lib/model';

describe('AuthReducer', () => {
  describe('undefined action', () => {
    it('should return the default state', () => {
      const action = {} as any;

      const result = reducer(undefined, action);

      /**
       * Snapshot tests are a quick way to validate
       * the state produced by a reducer since
       * its plain JavaScript object. These snapshots
       * are used to validate against the current state
       * if the functionality of the reducer ever changes.
       */
      expect(result).toMatchSnapshot();
    });
  });

  describe('LOGIN_SUCCESS', () => {
    it('should add a user set loggedIn to true in auth state', () => {
      const account = new Account();
      const createAction = new AuthApiActions.LoginSuccess({ account });

      const result = reducer(fromAuth.initialState, createAction);

      expect(result).toMatchSnapshot();
    });
  });

  describe('LOGOUT', () => {
    it('should logout a user', () => {
      const initialState = {
        account: new Account(),
      } as fromAuth.State;
      const createAction = new AuthActions.Logout();

      const result = reducer(initialState, createAction);

      expect(result).toMatchSnapshot();
    });
  });
});
