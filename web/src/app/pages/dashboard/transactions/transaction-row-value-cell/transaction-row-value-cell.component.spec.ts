import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TransactionRowValueCellComponent } from './transaction-row-value-cell.component';
import { CryptoService } from 'src/app/lib/services';
import { SharedModule } from 'src/app/lib/shared.module';
import { compileComponent } from '@angular/core/src/render3/jit/directive';
import { Message } from 'src/app/lib/model';
import { I18nService } from 'src/app/lib/i18n/i18n.service';
import { I18nPipe } from 'src/app/lib/i18n/i18n.pipe';
import { I18nModule } from 'src/app/lib/i18n/i18n.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('TransactionRowValueCellComponent', () => {
  let component: TransactionRowValueCellComponent;
  let fixture: ComponentFixture<TransactionRowValueCellComponent>;

  beforeEach(async(() => {
    
    TestBed.configureTestingModule({
      declarations: [ TransactionRowValueCellComponent ],
      imports: [ I18nModule, HttpClientTestingModule ],
      providers: [ CryptoService, I18nService ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TransactionRowValueCellComponent);
    component = fixture.componentInstance;
  });

  test('should create', () => {
    component.value = "when moon";
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  test('should render unencrypted messages', () => {
    component.value = new Message({
      message: ''
    });
    fixture.detectChanges();
    expect(component.valueType).toBe('Message');
  });
});
