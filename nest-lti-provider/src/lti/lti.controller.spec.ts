import { Test, TestingModule } from '@nestjs/testing';
import { LtiController } from './lti.controller';
import { LtiService } from './lti.service';

describe('LtiController', () => {
  let controller: LtiController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [LtiController],
      providers: [LtiService],
    }).compile();

    controller = module.get<LtiController>(LtiController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
