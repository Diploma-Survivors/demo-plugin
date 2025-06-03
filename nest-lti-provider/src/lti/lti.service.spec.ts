import { Test, TestingModule } from '@nestjs/testing';
import { LtiService } from './lti.service';

describe('LtiService', () => {
  let service: LtiService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [LtiService],
    }).compile();

    service = module.get<LtiService>(LtiService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
