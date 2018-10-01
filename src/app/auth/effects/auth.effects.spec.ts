import { TestBed } from '@angular/core/testing';
import { MatDialog } from '@angular/material';
import { Router } from '@angular/router';
import { Actions } from '@ngrx/effects';
import { provideMockActions } from '@ngrx/effects/testing';
import { cold, hot } from 'jest-marbles';
import { Observable, of } from 'rxjs';
import {
  LoginPageActions,
  AuthActions,
  AuthApiActions,
} from '../actions';

import { Credentials } from '../models/credentials';
import { AuthService } from '../services/auth.service';
import { AuthEffects } from '../effects/auth.effects';
import { Account } from '../../lib/model';

describe('AuthEffects', () => {
  let effects: AuthEffects;
  let authService: any;
  let actions$: Observable<any>;
  let routerService: any;
  let dialog: any;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        // AuthEffects,
        // {
        //   provide: AuthService,
        //   useValue: { login: jest.fn() },
        // },
        provideMockActions(() => actions$),
        {
          provide: Router,
          useValue: { navigate: jest.fn() },
        },
        {
          provide: MatDialog,
          useValue: {
            open: jest.fn(),
          }
        },
      ],
    });

    // effects = TestBed.get(AuthEffects);
    // authService = TestBed.get(AuthService);
    actions$ = TestBed.get(Actions);
    routerService = TestBed.get(Router);
    dialog = TestBed.get(MatDialog);

    spyOn(routerService, 'navigate').and.callThrough();
  });


  xdescribe('login$', () => {
    it('should return an auth.LoginSuccess action, with user information if login succeeds', () => {
      const credentials: Credentials = { passphrase: '' };
      const account = new Account();
      const action = new LoginPageActions.Login({ credentials });
      const completion = new AuthApiActions.LoginSuccess({ account });

      actions$ = hot('-a---', { a: action });
      const response = cold('-a|', { a: account });
      const expected = cold('--b', { b: completion });
      authService.login = jest.fn(() => response);

      expect(effects.login$).toBeObservable(expected);
    });

    xit('should return a new auth.LoginFailure if the login service throws', () => {
      const credentials: Credentials = { passphrase: '' };
      const action = new LoginPageActions.Login({ credentials });
      const completion = new AuthApiActions.LoginFailure({
        error: 'Invalid username or password',
      });
      const error = 'Invalid username or password';

      actions$ = hot('-a---', { a: action });
      const response = cold('-#', {}, error);
      const expected = cold('--b', { b: completion });
      authService.login = jest.fn(() => response);

      expect(effects.login$).toBeObservable(expected);
    });
  });

  xdescribe('loginSuccess$', () => {
    xit('should dispatch a RouterNavigation action', (done: any) => {
      const account = new Account();
      const action = new AuthApiActions.LoginSuccess({ account });

      actions$ = of(action);
      effects.loginSuccess$.subscribe(() => {
        expect(routerService.navigate).toHaveBeenCalledWith(['/']);
        done();
      });
    });
  });

  xdescribe('loginRedirect$', () => {
    xit('should dispatch a RouterNavigation action when auth.LoginRedirect is dispatched', (done: any) => {
      const action = new AuthApiActions.LoginRedirect();

      actions$ = of(action);

      effects.loginRedirect$.subscribe(() => {
        expect(routerService.navigate).toHaveBeenCalledWith(['/login']);
        done();
      });
    });

    xit('should dispatch a RouterNavigation action when auth.Logout is dispatched', (done: any) => {
      const action = new AuthActions.Logout();

      actions$ = of(action);

      effects.loginRedirect$.subscribe(() => {
        expect(routerService.navigate).toHaveBeenCalledWith(['/login']);
        done();
      });
    });
  });

  xdescribe('logoutConfirmation$', () => {

    
    xit('should dispatch a Logout action if dialog closes with true result', () => {
      const action = new AuthActions.LogoutConfirmation();
      const completion = new AuthActions.Logout();

      actions$ = hot('-a', { a: action });
      const expected = cold('-b', { b: completion });

      dialog.open = () => ({
        afterClosed: jest.fn(() => of(true)),
      });

      expect(effects.logoutConfirmation$).toBeObservable(expected);
    });

    xit('should dispatch a LogoutConfirmationDismiss action if dialog closes with falsy result', () => {
      const action = new AuthActions.LogoutConfirmation();
      const completion = new AuthActions.LogoutConfirmationDismiss();

      actions$ = hot('-a', { a: action });
      const expected = cold('-b', { b: completion });

      dialog.open = () => ({
        afterClosed: jest.fn(() => of(false)),
      });

      expect(effects.logoutConfirmation$).toBeObservable(expected);
    });
  });
});
